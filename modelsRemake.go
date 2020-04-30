package main

import "twitter/twitter/models"

func (tweet *Tweet) toModel() *models.Tweet {
	model := new(models.Tweet)
	model.ID = int64(tweet.ID)
	model.Text = tweet.Text
	model.AuthorID = int64(tweet.AuthorID)
	model.Author = tweet.Author
	return model
}

func tweetArrayToModel(tweets []*Tweet) []*models.Tweet {
	modelsArray := make([]*models.Tweet, len(tweets))
	for i := 0; i < len(tweets); i++ {
		modelsArray[i] = tweets[i].toModel()
	}
	return modelsArray
}
