package main

// var mockBatchChange = &BatchChange{}
var mockBatchChanges = []batchChange{
	{
		// public batch change created by Kai belonging to the Acme Corp Org namespace
		Name:           "kai-batchChange1",
		Private:        false,
		NamespaceOrgID: 1,
		CreatorID:      1,
	},
	{
		// private batch change created by Kai belonging to the Acme Corp Org namespace
		Name:           "kai-batchChange2",
		Private:        true,
		NamespaceOrgID: 1,
		CreatorID:      1,
	},
	// public batch change created by Kai
	{
		Name:            "kai-batchChange3",
		Private:         false,
		NamespaceUserID: 1,
		CreatorID:       1,
	},
	{
		Name:            "kai-batchChange4",
		Private:         false,
		NamespaceUserID: 1,
		CreatorID:       1,
	},
	// private batch change created by Kai
	{
		Name:            "kai-batchChange5",
		Private:         true,
		NamespaceUserID: 1,
		CreatorID:       1,
	},
	// private batch change created by elliot
	{
		Name:            "elliot-batchChange1",
		Private:         true,
		NamespaceUserID: 2,
		CreatorID:       2,
	},
	// public batch change created by elliot
	{
		Name:            "elliot-batchChange2",
		Private:         false,
		NamespaceUserID: 2,
		CreatorID:       2,
	},
}

var mockCodeInsights = []codeinsight{
	{
		Name:   "kai-codeinsights1",
		UserID: 1,
	},
	{
		Name:   "kai-codeinsights2",
		UserID: 1,
	},
}

var mockNotebooks = []notebook{
	{
		Name:      "kai-notebook1",
		Content:   "kai-notebook1-content",
		Private:   false,
		CreatorID: 1,
	},
	{
		Name:      "elliot-notebook1",
		Content:   "elliot-notebook1-content",
		Private:   true,
		CreatorID: 2,
	},
	{
		Name:      "elliot-notebook2",
		Content:   "elliot-notebook2-content",
		Private:   true,
		CreatorID: 2,
	},
	{
		Name:      "elliot-notebook3",
		Content:   "elliot-notebook3-content",
		Private:   true,
		CreatorID: 2,
	},
	{
		Name:      "jalen-notebook1",
		Content:   "jalen-notebook1-content",
		Private:   false,
		CreatorID: 3,
	},
}
