package kv

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/influxdata/influxdb"
	icontext "github.com/influxdata/influxdb/context"
	"github.com/influxdata/influxdb/kit/tracing"
	"go.uber.org/zap"
)

var (
	urmBucket      = []byte("userresourcemappingsv1")
	urmIndexBucket = []byte("userresourcemappingsindexv1")

	// ErrInvalidURMID is used when the service was provided
	// an invalid ID format.
	ErrInvalidURMID = &influxdb.Error{
		Code: influxdb.EInvalid,
		Msg:  "provided user resource mapping ID has invalid format",
	}

	// ErrURMNotFound is used when the user resource mapping is not found.
	ErrURMNotFound = &influxdb.Error{
		Msg:  "user to resource mapping not found",
		Code: influxdb.ENotFound,
	}
)

// UnavailableURMServiceError is used if we aren't able to interact with the
// store, it means the store is not available at the moment (e.g. network).
func UnavailableURMServiceError(err error) *influxdb.Error {
	return &influxdb.Error{
		Code: influxdb.EInternal,
		Msg:  fmt.Sprintf("Unable to connect to resource mapping service. Please try again; Err: %v", err),
		Op:   "kv/userResourceMapping",
	}
}

// CorruptURMError is used when the config cannot be unmarshalled from the
// bytes stored in the kv.
func CorruptURMError(err error) *influxdb.Error {
	return &influxdb.Error{
		Code: influxdb.EInternal,
		Msg:  fmt.Sprintf("Unknown internal user resource mapping data error; Err: %v", err),
		Op:   "kv/userResourceMapping",
	}
}

// ErrUnprocessableMapping is used when a user resource mapping  is not able to be converted to JSON.
func ErrUnprocessableMapping(err error) *influxdb.Error {
	return &influxdb.Error{
		Code: influxdb.EUnprocessableEntity,
		Msg:  fmt.Sprintf("unable to convert mapping of user to resource into JSON; Err %v", err),
	}
}

// NonUniqueMappingError is an internal error when a user already has
// been mapped to a resource
func NonUniqueMappingError(userID influxdb.ID) error {
	return &influxdb.Error{
		Code: influxdb.EInternal,
		Msg:  fmt.Sprintf("Unexpected error when assigning user to a resource: mapping for user %s already exists", userID.String()),
	}
}

func (s *Service) initializeURMs(ctx context.Context, tx Tx) error {
	if _, err := tx.Bucket(urmBucket); err != nil {
		return UnavailableURMServiceError(err)
	}
	if _, err := tx.Bucket(urmIndexBucket); err != nil {
		return UnavailableURMServiceError(err)
	}
	return nil
}

func filterMappingsFn(filter influxdb.UserResourceMappingFilter) func(m *influxdb.UserResourceMapping) bool {
	return func(mapping *influxdb.UserResourceMapping) bool {
		return (!filter.UserID.Valid() || (filter.UserID == mapping.UserID)) &&
			(!filter.ResourceID.Valid() || (filter.ResourceID == mapping.ResourceID)) &&
			(filter.UserType == "" || (filter.UserType == mapping.UserType)) &&
			(filter.ResourceType == "" || (filter.ResourceType == mapping.ResourceType))
	}
}

// FindUserResourceMappings returns a list of UserResourceMappings that match filter and the total count of matching mappings.
func (s *Service) FindUserResourceMappings(ctx context.Context, filter influxdb.UserResourceMappingFilter, opt ...influxdb.FindOptions) ([]*influxdb.UserResourceMapping, int, error) {
	var ms []*influxdb.UserResourceMapping
	err := s.kv.View(ctx, func(tx Tx) error {
		var err error
		ms, err = s.findUserResourceMappings(ctx, tx, filter)
		return err
	})

	if err != nil {
		return nil, 0, err
	}

	return ms, len(ms), nil
}

func userResourceMappingPredicate(filter influxdb.UserResourceMappingFilter) CursorPredicateFunc {
	switch {
	case filter.ResourceID.Valid() && filter.UserID.Valid():
		keyPredicate := filter.ResourceID.String() + filter.UserID.String()
		return func(key, _ []byte) bool {
			return len(key) >= 32 && string(key[:32]) == keyPredicate
		}

	case !filter.ResourceID.Valid() && filter.UserID.Valid():
		keyPredicate := filter.UserID.String()
		return func(key, _ []byte) bool {
			return len(key) >= 32 && string(key[16:32]) == keyPredicate
		}

	case filter.ResourceID.Valid() && !filter.UserID.Valid():
		keyPredicate := filter.ResourceID.String()
		return func(key, _ []byte) bool {
			return len(key) >= 16 && string(key[:16]) == keyPredicate
		}

	default:
		return nil
	}
}

type urmFindOptions struct {
	skipKeys map[string]struct{}
}

type urmFindOption func(*urmFindOptions)

func (f *urmFindOptions) skip(key string) (skip bool) {
	if f.skipKeys == nil {
		return false
	}

	_, skip = f.skipKeys[key]
	return
}

func withSkipKey(key string) urmFindOption {
	return func(o *urmFindOptions) {
		if o.skipKeys == nil {
			o.skipKeys = map[string]struct{}{}
		}

		o.skipKeys[key] = struct{}{}
	}
}

func (s *Service) findUserResourceMappings(ctx context.Context, tx Tx, filter influxdb.UserResourceMappingFilter, opts ...urmFindOption) ([]*influxdb.UserResourceMapping, error) {
	ms := []*influxdb.UserResourceMapping{}
	pred := userResourceMappingPredicate(filter)
	filterFn := filterMappingsFn(filter)
	// if we are given a user id we should try finding by index
	if filter.UserID.Valid() {
		var err error
		ms, err = s.findUserResourceMappingsByIndex(ctx, tx, filter, opts...)
		if err != nil {
			return nil, err
		}

		// if we found nothing we need to fall back on the old method because the index may not have been created
		if len(ms) > 0 {
			return ms, nil
		}
	}

	err := s.forEachUserResourceMapping(ctx, tx, pred, func(m *influxdb.UserResourceMapping) bool {
		if filterFn(m) {
			ms = append(ms, m)
		}
		return true
	})

	// if we got to this point we failed to find the user by the index so we need to populate the index
	if filter.UserID.Valid() && len(ms) > 0 {
		indexes := map[string][]byte{}
		for _, m := range ms {
			key, _ := userResourceKey(m)
			ikey, _ := userResourceIndexKey(m)
			indexes[string(ikey)] = key

		}

		s.indexer.AddToIndex(urmIndexBucket, indexes)
	}
	return ms, err
}

func (s *Service) findUserResourceMappingsByIndex(ctx context.Context, tx Tx, filter influxdb.UserResourceMappingFilter, opts ...urmFindOption) ([]*influxdb.UserResourceMapping, error) {
	var (
		ms       = []*influxdb.UserResourceMapping{}
		filterFn = filterMappingsFn(filter)
		options  = urmFindOptions{}
	)

	for _, opt := range opts {
		opt(&options)
	}

	bkt, err := tx.Bucket(urmBucket)
	if err != nil {
		return nil, err
	}

	idx, err := tx.Bucket(urmIndexBucket)
	if err != nil {
		return nil, err
	}

	prefix := urmIndexPrefix(filter.UserID)
	wrapInternal := func(err error) *influxdb.Error {
		return &influxdb.Error{
			Code: influxdb.EInternal,
			Err:  err,
		}
	}

	// index scan
	cursor, err := idx.ForwardCursor(prefix, WithCursorPrefix(prefix))
	if err != nil {
		return nil, wrapInternal(err)
	}

	for k, v := cursor.Next(); k != nil && v != nil; k, v = cursor.Next() {
		// step over skip keys
		if options.skip(string(v)) {
			continue
		}

		nv, err := bkt.Get(v)
		if err != nil {
			s.log.Info(
				"key not found",
				zap.String("function", "findUserResourceMappingsByIndex"),
				zap.String("key", string(v)),
			)
			continue
		}

		m := &influxdb.UserResourceMapping{}
		if err := json.Unmarshal(nv, m); err != nil {
			return nil, CorruptURMError(err)
		}

		if filterFn(m) {
			ms = append(ms, m)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, wrapInternal(err)
	}

	if err := cursor.Close(); err != nil {
		return nil, wrapInternal(err)
	}

	return ms, nil
}

func (s *Service) findUserResourceMapping(ctx context.Context, tx Tx, filter influxdb.UserResourceMappingFilter) (*influxdb.UserResourceMapping, error) {
	ms, err := s.findUserResourceMappings(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	if len(ms) == 0 {
		return nil, ErrURMNotFound
	}

	return ms[0], nil
}

// CreateUserResourceMapping associates a user to a resource either as a member
// or owner.
func (s *Service) CreateUserResourceMapping(ctx context.Context, m *influxdb.UserResourceMapping) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		return s.createUserResourceMapping(ctx, tx, m)
	})
}

func (s *Service) createUserResourceMapping(ctx context.Context, tx Tx, m *influxdb.UserResourceMapping) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if err := s.uniqueUserResourceMapping(ctx, tx, m); err != nil {
		return err
	}

	v, err := json.Marshal(m)
	if err != nil {
		return ErrUnprocessableMapping(err)
	}

	key, err := userResourceKey(m)
	if err != nil {
		return err
	}

	b, err := tx.Bucket(urmBucket)
	if err != nil {
		return UnavailableURMServiceError(err)
	}

	if err := b.Put(key, v); err != nil {
		return UnavailableURMServiceError(err)
	}

	ikey, err := userResourceIndexKey(m)
	if err != nil {
		return err
	}

	ib, err := tx.Bucket(urmIndexBucket)
	if err != nil {
		return UnavailableURMServiceError(err)
	}

	if err := ib.Put(ikey, key); err != nil {
		return UnavailableURMServiceError(err)
	}

	if m.ResourceType == influxdb.OrgsResourceType {
		return s.createOrgDependentMappings(ctx, tx, m)
	}

	return nil
}

// This method creates the user/resource mappings for resources that belong to an organization.
func (s *Service) createOrgDependentMappings(ctx context.Context, tx Tx, m *influxdb.UserResourceMapping) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	bf := influxdb.BucketFilter{OrganizationID: &m.ResourceID}
	bs, err := s.findBuckets(ctx, tx, bf)
	if err != nil {
		return err
	}
	for _, b := range bs {
		m := &influxdb.UserResourceMapping{
			ResourceType: influxdb.BucketsResourceType,
			ResourceID:   b.ID,
			UserType:     m.UserType,
			UserID:       m.UserID,
		}
		if err := s.createUserResourceMapping(ctx, tx, m); err != nil {
			return err
		}
		// TODO(desa): add support for all other resource types.
	}

	return nil
}

func userResourceKey(m *influxdb.UserResourceMapping) ([]byte, error) {
	encodedResourceID, err := m.ResourceID.Encode()
	if err != nil {
		return nil, ErrInvalidURMID
	}

	encodedUserID, err := m.UserID.Encode()
	if err != nil {
		return nil, ErrInvalidURMID
	}

	key := make([]byte, len(encodedResourceID)+len(encodedUserID))
	copy(key, encodedResourceID)
	copy(key[len(encodedResourceID):], encodedUserID)

	return key, nil
}

func userResourceIndexKey(m *influxdb.UserResourceMapping) ([]byte, error) {
	encodedResourceID, err := m.ResourceID.Encode()
	if err != nil {
		return nil, ErrInvalidURMID
	}

	encodedUserID, err := m.UserID.Encode()
	if err != nil {
		return nil, ErrInvalidURMID
	}

	key := append(encodedUserID, '/')
	return append(key, encodedResourceID...), nil
}

func urmIndexPrefix(userID influxdb.ID) []byte {
	id, _ := userID.Encode()
	return append(id, '/')
}

func (s *Service) forEachUserResourceMapping(ctx context.Context, tx Tx, pred CursorPredicateFunc, fn func(*influxdb.UserResourceMapping) bool) error {
	b, err := tx.Bucket(urmBucket)
	if err != nil {
		return UnavailableURMServiceError(err)
	}
	var cur Cursor
	if pred != nil {
		cur, err = b.Cursor(WithCursorHintPredicate(pred))
	} else {
		cur, err = b.Cursor()
	}
	if err != nil {
		return UnavailableURMServiceError(err)
	}

	for k, v := cur.First(); k != nil; k, v = cur.Next() {
		m := &influxdb.UserResourceMapping{}
		if err := json.Unmarshal(v, m); err != nil {
			return CorruptURMError(err)
		}

		if !fn(m) {
			break
		}
	}

	return nil
}

func (s *Service) uniqueUserResourceMapping(ctx context.Context, tx Tx, m *influxdb.UserResourceMapping) error {
	key, err := userResourceKey(m)
	if err != nil {
		return err
	}

	b, err := tx.Bucket(urmBucket)
	if err != nil {
		return UnavailableURMServiceError(err)
	}

	_, err = b.Get(key)
	if !IsNotFound(err) {
		return NonUniqueMappingError(m.UserID)
	}

	return nil
}

// DeleteUserResourceMapping deletes a user resource mapping.
func (s *Service) DeleteUserResourceMapping(ctx context.Context, resourceID influxdb.ID, userID influxdb.ID) error {
	return s.kv.Update(ctx, func(tx Tx) error {
		// TODO(goller): I don't think this find is needed as delete also finds.
		m, err := s.findUserResourceMapping(ctx, tx, influxdb.UserResourceMappingFilter{
			ResourceID: resourceID,
			UserID:     userID,
		})
		if err != nil {
			return err
		}

		filter := influxdb.UserResourceMappingFilter{
			ResourceID: resourceID,
			UserID:     userID,
		}
		if err := s.deleteUserResourceMapping(ctx, tx, filter); err != nil {
			return err
		}

		if m.ResourceType == influxdb.OrgsResourceType {
			key, err := userResourceKey(m)
			if err != nil {
				// I'm not super concerned that we will get here.  We know this is a valid resource
				// because we've just found it above.  Me of the future... if this was a problem,
				// sorry.
				return err
			}
			return s.deleteOrgDependentMappings(ctx, tx, m, withSkipKey(string(key)))
		}

		return nil
	})
}

func (s *Service) deleteUserResourceMapping(ctx context.Context, tx Tx, filter influxdb.UserResourceMappingFilter, opts ...urmFindOption) error {
	// TODO(goller): do we really need to find here? Seems like a Get is
	// good enough.
	ms, err := s.findUserResourceMappings(ctx, tx, filter, opts...)
	if err != nil {
		return err
	}
	if len(ms) == 0 {
		return ErrURMNotFound
	}

	key, err := userResourceKey(ms[0])
	if err != nil {
		return err
	}

	ikey, err := userResourceIndexKey(ms[0])
	if err != nil {
		return err
	}

	b, err := tx.Bucket(urmBucket)
	if err != nil {
		return UnavailableURMServiceError(err)
	}

	ib, err := tx.Bucket(urmIndexBucket)
	if err != nil {
		return UnavailableURMServiceError(err)
	}
	_, err = b.Get(key)
	if IsNotFound(err) {
		return ErrURMNotFound
	}
	if err != nil {
		return UnavailableURMServiceError(err)
	}

	if err := b.Delete(key); err != nil {
		return UnavailableURMServiceError(err)
	}

	if err := ib.Delete(ikey); err != nil {
		return UnavailableURMServiceError(err)
	}

	return nil
}

func (s *Service) deleteUserResourceMappings(ctx context.Context, tx Tx, filter influxdb.UserResourceMappingFilter) error {
	ms, err := s.findUserResourceMappings(ctx, tx, filter)
	if err != nil {
		return err
	}
	for _, m := range ms {
		key, err := userResourceKey(m)
		if err != nil {
			return err
		}

		ikey, err := userResourceIndexKey(m)
		if err != nil {
			return err
		}

		b, err := tx.Bucket(urmBucket)
		if err != nil {
			return UnavailableURMServiceError(err)
		}

		_, err = b.Get(key)
		if IsNotFound(err) {
			return ErrURMNotFound
		}
		if err != nil {
			return UnavailableURMServiceError(err)
		}

		ib, err := tx.Bucket(urmIndexBucket)
		if err != nil {
			return UnavailableURMServiceError(err)
		}

		if err := b.Delete(key); err != nil {
			return UnavailableURMServiceError(err)
		}
		if err := ib.Delete(ikey); err != nil {
			return UnavailableURMServiceError(err)
		}
	}
	return nil
}

// This method deletes the user/resource mappings for resources that belong to an organization.
func (s *Service) deleteOrgDependentMappings(ctx context.Context, tx Tx, m *influxdb.UserResourceMapping, opts ...urmFindOption) error {
	bf := influxdb.BucketFilter{OrganizationID: &m.ResourceID}
	bs, err := s.findBuckets(ctx, tx, bf)
	if err != nil {
		return err
	}
	for _, b := range bs {
		filter := influxdb.UserResourceMappingFilter{
			ResourceType: influxdb.BucketsResourceType,
			ResourceID:   b.ID,
			UserID:       m.UserID,
		}

		if err := s.deleteUserResourceMapping(ctx, tx, filter, opts...); err != nil {
			if influxdb.ErrorCode(err) == influxdb.ENotFound {
				s.log.Info("URM bucket is missing", zap.Stringer("orgID", m.ResourceID))
				continue
			}
			return err
		}
		// TODO(desa): add support for all other resource types.
	}

	return nil
}

func (s *Service) addResourceOwner(ctx context.Context, tx Tx, rt influxdb.ResourceType, id influxdb.ID) error {
	a, err := icontext.GetAuthorizer(ctx)
	if err != nil {
		return &influxdb.Error{
			Code: influxdb.EInternal,
			Msg:  fmt.Sprintf("could not find authorizer on context when adding user to resource type %s", rt),
		}
	}

	urm := &influxdb.UserResourceMapping{
		ResourceID:   id,
		ResourceType: rt,
		UserID:       a.GetUserID(),
		UserType:     influxdb.Owner,
	}

	if err := s.createUserResourceMapping(ctx, tx, urm); err != nil {
		return &influxdb.Error{
			Code: influxdb.EInternal,
			Msg:  "could not create user resource mapping",
			Err:  err,
		}
	}

	return nil
}
