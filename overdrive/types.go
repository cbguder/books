package overdrive

import "time"

type MediaResponse struct {
	Items []struct {
		IsAvailable              bool `json:"isAvailable"`
		IsPreReleaseTitle        bool `json:"isPreReleaseTitle"`
		IsFastlane               bool `json:"isFastlane"`
		IsRecommendableToLibrary bool `json:"isRecommendableToLibrary"`
		IsOwned                  bool `json:"isOwned"`
		IsHoldable               bool `json:"isHoldable"`
		IsAdvantageFiltered      bool `json:"isAdvantageFiltered"`
		VisitorEligible          bool `json:"visitorEligible"`
		JuvenileEligible         bool `json:"juvenileEligible"`
		YoungAdultEligible       bool `json:"youngAdultEligible"`
		FirstCreatorId           int  `json:"firstCreatorId"`
		IsBundledChild           bool `json:"isBundledChild"`
		Subjects                 []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"subjects"`
		BisacCodes []string `json:"bisacCodes"`
		Bisac      []struct {
			Description string `json:"description"`
			Code        string `json:"code"`
		} `json:"bisac"`
		Levels []struct {
			Id    string `json:"id"`
			Name  string `json:"name"`
			Value string `json:"value"`
			High  string `json:"high,omitempty"`
			Low   string `json:"low,omitempty"`
		} `json:"levels"`
		Creators []struct {
			Id       int    `json:"id"`
			SortName string `json:"sortName"`
			Role     string `json:"role"`
			Name     string `json:"name"`
		} `json:"creators"`
		Languages []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"languages"`
		Imprint struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"imprint"`
		Ratings struct {
			MaturityLevel struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"maturityLevel"`
			NaughtyScore struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"naughtyScore"`
		} `json:"ratings"`
		ReviewCounts struct {
			PublisherSupplier int `json:"publisherSupplier"`
			Premium           int `json:"premium"`
		} `json:"reviewCounts"`
		Awards []struct {
			Id          int    `json:"id"`
			Source      string `json:"source"`
			Description string `json:"description"`
		} `json:"awards"`
		Constraints struct {
			IsDisneyEulaRequired bool `json:"isDisneyEulaRequired"`
		} `json:"constraints"`
		Sample struct {
			Href string `json:"href"`
		} `json:"sample"`
		Publisher struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"publisher"`
		Type struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"type"`
		Covers struct {
			Cover150Wide struct {
				PrimaryColor struct {
					Rgb struct {
						Blue  int `json:"blue"`
						Green int `json:"green"`
						Red   int `json:"red"`
					} `json:"rgb"`
					Hex string `json:"hex"`
				} `json:"primaryColor"`
				IsPlaceholderImage bool   `json:"isPlaceholderImage"`
				Width              int    `json:"width"`
				Height             int    `json:"height"`
				Href               string `json:"href"`
			} `json:"cover150Wide"`
			Cover300Wide struct {
				PrimaryColor struct {
					Rgb struct {
						Blue  int `json:"blue"`
						Green int `json:"green"`
						Red   int `json:"red"`
					} `json:"rgb"`
					Hex string `json:"hex"`
				} `json:"primaryColor"`
				IsPlaceholderImage bool   `json:"isPlaceholderImage"`
				Width              int    `json:"width"`
				Height             int    `json:"height"`
				Href               string `json:"href"`
			} `json:"cover300Wide"`
			Cover510Wide struct {
				PrimaryColor struct {
					Rgb struct {
						Blue  int `json:"blue"`
						Green int `json:"green"`
						Red   int `json:"red"`
					} `json:"rgb"`
					Hex string `json:"hex"`
				} `json:"primaryColor"`
				IsPlaceholderImage bool   `json:"isPlaceholderImage"`
				Width              int    `json:"width"`
				Height             int    `json:"height"`
				Href               string `json:"href"`
			} `json:"cover510Wide"`
		} `json:"covers"`
		Formats []struct {
			IsBundleParent           bool `json:"isBundleParent"`
			HasAudioSynchronizedText bool `json:"hasAudioSynchronizedText"`
			Identifiers              []struct {
				Value string `json:"value"`
				Type  string `json:"type"`
			} `json:"identifiers"`
			Rights          []interface{} `json:"rights"`
			BundledContent  []interface{} `json:"bundledContent"`
			Id              string        `json:"id"`
			Name            string        `json:"name"`
			OnSaleDateUtc   time.Time     `json:"onSaleDateUtc"`
			FulfillmentType string        `json:"fulfillmentType"`
			Sample          struct {
				Href string `json:"href"`
			} `json:"sample,omitempty"`
			Isbn     string `json:"isbn,omitempty"`
			FileSize int    `json:"fileSize,omitempty"`
		} `json:"formats"`
		PublisherAccount struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"publisherAccount"`
		EstimatedReleaseDate    time.Time `json:"estimatedReleaseDate"`
		Subtitle                string    `json:"subtitle"`
		Description             string    `json:"description"`
		AvailableCopies         int       `json:"availableCopies"`
		OwnedCopies             int       `json:"ownedCopies"`
		LuckyDayAvailableCopies int       `json:"luckyDayAvailableCopies"`
		LuckyDayOwnedCopies     int       `json:"luckyDayOwnedCopies"`
		HoldsCount              int       `json:"holdsCount"`
		HoldsRatio              int       `json:"holdsRatio"`
		EstimatedWaitDays       int       `json:"estimatedWaitDays"`
		AvailabilityType        string    `json:"availabilityType"`
		Id                      string    `json:"id"`
		FirstCreatorName        string    `json:"firstCreatorName"`
		FirstCreatorSortName    string    `json:"firstCreatorSortName"`
		Title                   string    `json:"title"`
		SortTitle               string    `json:"sortTitle"`
		StarRating              float64   `json:"starRating"`
		StarRatingCount         int       `json:"starRatingCount"`
		PublishDate             time.Time `json:"publishDate"`
		PublishDateText         string    `json:"publishDateText"`
		ReserveId               string    `json:"reserveId"`
	} `json:"items"`
	Links struct {
		Self struct {
			Page     int    `json:"page"`
			PageText string `json:"pageText"`
		} `json:"self"`
		First struct {
			Page     int    `json:"page"`
			PageText string `json:"pageText"`
		} `json:"first"`
		Last struct {
			Page     int    `json:"page"`
			PageText string `json:"pageText"`
		} `json:"last"`
	} `json:"links"`
}
