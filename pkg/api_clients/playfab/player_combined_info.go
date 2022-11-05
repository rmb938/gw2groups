package playfab

type PlayerProfileViewConstraintsRequest struct {
	ShowAvatarUrl                     *bool `json:"ShowAvatarUrl,omitempty"`
	ShowBannedUntil                   *bool `json:"ShowBannedUntil,omitempty"`
	ShowCampaignAttributions          *bool `json:"ShowCampaignAttributions,omitempty"`
	ShowContactEmailAddresses         *bool `json:"ShowContactEmailAddresses,omitempty"`
	ShowCreated                       *bool `json:"ShowCreated,omitempty"`
	ShowDisplayName                   *bool `json:"ShowDisplayName,omitempty"`
	ShowExperimentVariants            *bool `json:"ShowExperimentVariants,omitempty"`
	ShowLastLogin                     *bool `json:"ShowLastLogin,omitempty"`
	ShowLinkedAccounts                *bool `json:"ShowLinkedAccounts,omitempty"`
	ShowLocations                     *bool `json:"ShowLocations,omitempty"`
	ShowMemberships                   *bool `json:"ShowMemberships,omitempty"`
	ShowOrigination                   *bool `json:"ShowOrigination,omitempty"`
	ShowPushNotificationRegistrations *bool `json:"ShowPushNotificationRegistrations,omitempty"`
	ShowStatistics                    *bool `json:"ShowStatistics,omitempty"`
	ShowTags                          *bool `json:"ShowTags,omitempty"`
	ShowTotalValueToDateInUsd         *bool `json:"ShowTotalValueToDateInUsd,omitempty"`
	ShowValuesToDate                  *bool `json:"ShowValuesToDate,omitempty"`
}

type PlayerCombinedInfoRequestParams struct {
	GetCharacterInventories *bool                                `json:"GetCharacterInventories,omitempty"`
	GetCharacterList        *bool                                `json:"GetCharacterList,omitempty"`
	GetPlayerProfile        *bool                                `json:"GetPlayerProfile,omitempty"`
	GetPlayerStatistics     *bool                                `json:"GetPlayerStatistics,omitempty"`
	GetTitleData            *bool                                `json:"GetTitleData,omitempty"`
	GetUserAccountInfo      *bool                                `json:"GetUserAccountInfo,omitempty"`
	GetUserData             *bool                                `json:"GetUserData,omitempty"`
	GetUserInventory        *bool                                `json:"GetUserInventory,omitempty"`
	GetUserReadOnlyData     *bool                                `json:"GetUserReadOnlyData,omitempty"`
	GetUserVirtualCurrency  *bool                                `json:"GetUserVirtualCurrency,omitempty"`
	PlayerStatisticNames    []string                             `json:"PlayerStatisticNames,omitempty"`
	ProfileConstraints      *PlayerProfileViewConstraintsRequest `json:"ProfileConstraints,omitempty"`
	TitleDataKeys           []string                             `json:"TitleDataKeys,omitempty"`
	UserDataKeys            []string                             `json:"UserDataKeys,omitempty"`
	UserReadOnlyDataKeys    []string                             `json:"UserReadOnlyDataKeys,omitempty"`
}

type PlayerCombinedInfoResponse struct {
	AccountInfo                      interface{}                         `json:"AccountInfo"`
	CharacterInventories             []interface{}                       `json:"CharacterInventories"`
	CharacterList                    []interface{}                       `json:"CharacterList"`
	PlayerProfile                    interface{}                         `json:"PlayerProfile"`
	PlayerStatistics                 []PlayerStatisticValueResponse      `json:"PlayerStatistics"`
	TitleData                        map[string]string                   `json:"TitleData"`
	UserData                         map[string]UserDataRecordResponse   `json:"UserData"`
	UserDataVersion                  int                                 `json:"UserDataVersion"`
	UserInventory                    []ItemInstanceResponse              `json:"UserInventory"`
	UserReadOnlyData                 map[string]UserDataRecordResponse   `json:"UserReadOnlyData"`
	UserReadOnlyDataVersion          int                                 `json:"UserReadOnlyDataVersion"`
	UserVirtualCurrency              map[string]int                      `json:"UserVirtualCurrency"`
	UserVirtualCurrencyRechargeTimes VirtualCurrencyRechargeTimeResponse `json:"UserVirtualCurrencyRechargeTimes"`
}
