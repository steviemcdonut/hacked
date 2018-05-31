package resource_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ResourceChangeNotifierSuite struct {
	suite.Suite

	lastModifiedIDs []resource.ID
	lastFailedIDs   []resource.ID

	localizedResources resource.LocalizedResourcesList

	instance resource.ResourceChangeNotifier
}

func TestResourceChangeNotifierSuite(t *testing.T) {
	suite.Run(t, new(ResourceChangeNotifierSuite))
}

func (suite *ResourceChangeNotifierSuite) SetupTest() {
	suite.lastModifiedIDs = nil
	suite.lastFailedIDs = nil

	suite.instance.Localizer = suite
	suite.instance.Callback = suite.onResourcesModified
}

func (suite *ResourceChangeNotifierSuite) TestNoIDsAreNotifiedIfDataNotChanged() {
	suite.givenListIsMadeOf(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, [][]byte{{0x01}})))
	suite.whenResourcesAreModified(func() {}, 0x1000)
	suite.thenModifiedIDsShouldBe()
}

func (suite *ResourceChangeNotifierSuite) TestIDsAreNotifiedOfNewData() {
	suite.givenListIsMadeOf(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, [][]byte{{0x01}})))
	suite.whenResourcesAreModified(func() {
		suite.extendingListWith(
			suite.someLocalizedResources(resource.LangGerman,
				suite.storing(0x2000, [][]byte{{0x01}})))
	}, 0x2000)
	suite.thenModifiedIDsShouldBe(0x2000)
}

func (suite *ResourceChangeNotifierSuite) TestIDsAreNotifiedOfRemovedData() {
	suite.givenListIsMadeOf(
		suite.someLocalizedResources(resource.LangGerman,
			suite.storing(0x1000, [][]byte{{0x01}})))
	suite.whenResourcesAreModified(func() {
		suite.listIsMadeOf()
	}, 0x1000)
	suite.thenModifiedIDsShouldBe(0x1000)
}

func (suite *ResourceChangeNotifierSuite) TestIDsAreNotifiedOfChangedData() {
	suite.givenListIsMadeOf(
		suite.someLocalizedResources(resource.LangAny,
			suite.storing(0x1000, [][]byte{{0x01}})))
	suite.whenResourcesAreModified(func() {
		suite.listIsMadeOf(
			suite.someLocalizedResources(resource.LangAny,
				suite.storing(0x1000, [][]byte{{0x02}})))
	}, 0x1000)
	suite.thenModifiedIDsShouldBe(0x1000)
}

func (suite *ResourceChangeNotifierSuite) givenListIsMadeOf(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *ResourceChangeNotifierSuite) extendingListWith(resources ...resource.LocalizedResources) {
	suite.localizedResources = append(suite.localizedResources, resources...)
}

func (suite *ResourceChangeNotifierSuite) listIsMadeOf(resources ...resource.LocalizedResources) {
	suite.localizedResources = resources
}

func (suite *ResourceChangeNotifierSuite) whenResourcesAreModified(modifier func(), affectedIDs ...int) {
	convertedIDs := make([]resource.ID, len(affectedIDs))
	for i, id := range affectedIDs {
		convertedIDs[i] = resource.ID(id)
	}
	suite.instance.ModifyAndNotify(modifier, convertedIDs)
}

func (suite *ResourceChangeNotifierSuite) thenModifiedIDsShouldBe(expected ...int) {
	convertedIDs := make([]resource.ID, len(expected))
	for i, id := range expected {
		convertedIDs[i] = resource.ID(id)
	}

	assert.Equal(suite.T(), convertedIDs, suite.lastModifiedIDs, "IDs don't match")
}

func (suite *ResourceChangeNotifierSuite) someLocalizedResources(lang resource.Language, modifiers ...func(*resource.Store)) resource.LocalizedResources {
	store := resource.NewProviderBackedStore(resource.NullProvider())
	for _, modifier := range modifiers {
		modifier(store)
	}
	return resource.LocalizedResources{
		ID:       "unnamed",
		Language: lang,
		Provider: store,
	}
}

func (suite *ResourceChangeNotifierSuite) storing(id int, blocks [][]byte) func(*resource.Store) {
	return func(store *resource.Store) {
		store.Put(resource.ID(id), &resource.Resource{
			BlockProvider: resource.MemoryBlockProvider(blocks),
		})
	}
}

func (suite *ResourceChangeNotifierSuite) LocalizedResources(lang resource.Language) resource.ResourceSelector {
	return resource.ResourceSelector{
		From: suite.localizedResources,
		Lang: lang,
	}
}

func (suite *ResourceChangeNotifierSuite) onResourcesModified(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	suite.lastModifiedIDs = modifiedIDs
	suite.lastFailedIDs = failedIDs
}
