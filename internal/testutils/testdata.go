package testutils

import (
	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/store"
)

// TestData provides common test data fixtures for unit tests
type TestData struct{}

// NewTestData creates a new TestData instance
func NewTestData() *TestData {
	return &TestData{}
}

// CreateSampleProperty creates a sample property for testing
func (td *TestData) CreateSampleProperty() *cupid.Property {
	return &cupid.Property{
		HotelID:     12345,
		CupidID:     12345,
		HotelName:   "Test Hotel London",
		HotelType:   "Hotels",
		Chain:       "Test Chain",
		Latitude:    51.5074,
		Longitude:   -0.1278,
		Stars:       5,
		Rating:      9.5,
		ReviewCount: 100,
		Address: cupid.Address{
			Address:    "123 Test Street",
			City:       "London",
			State:      "England",
			Country:    "gb",
			PostalCode: "SW1A 1AA",
		},
		MainImageTh: "https://example.com/image.jpg",
	}
}

// CreateSamplePropertyData creates a sample property data with reviews and translations
func (td *TestData) CreateSamplePropertyData() *cupid.PropertyData {
	return &cupid.PropertyData{
		Property: *td.CreateSampleProperty(),
		Reviews: []cupid.Review{
			{
				ReviewID:     1,
				AverageScore: 9,
				Country:      "GB",
				Name:         "John Doe",
				Headline:     "Great hotel!",
				Pros:         "Clean, comfortable",
				Cons:         "No complaints",
				Date:         "2024-01-15",
				Language:     "en",
			},
			{
				ReviewID:     2,
				AverageScore: 8,
				Country:      "US",
				Name:         "Jane Smith",
				Headline:     "Good experience",
				Pros:         "Nice location",
				Cons:         "Could be better",
				Date:         "2024-01-10",
				Language:     "en",
			},
		},
		Translations: map[string]*cupid.Property{
			"fr": {
				HotelID:   12345,
				HotelName: "Hôtel de Test Londres",
				Address: cupid.Address{
					City:    "Londres",
					Country: "gb",
				},
			},
			"es": {
				HotelID:   12345,
				HotelName: "Hotel de Prueba Londres",
				Address: cupid.Address{
					City:    "Londres",
					Country: "gb",
				},
			},
		},
	}
}

// CreateSampleReview creates a sample review for testing
func (td *TestData) CreateSampleReview() cupid.Review {
	return cupid.Review{
		ReviewID:     1,
		AverageScore: 9,
		Country:      "GB",
		Name:         "John Doe",
		Headline:     "Great hotel!",
		Pros:         "Clean, comfortable",
		Cons:         "No complaints",
		Date:         "2024-01-15",
		Language:     "en",
	}
}

// CreateSampleReviews creates multiple sample reviews for testing
func (td *TestData) CreateSampleReviews() []cupid.Review {
	return []cupid.Review{
		{
			ReviewID:     1,
			AverageScore: 9,
			Country:      "GB",
			Name:         "John Doe",
			Headline:     "Great hotel!",
			Pros:         "Clean, comfortable",
			Cons:         "No complaints",
			Date:         "2024-01-15",
			Language:     "en",
		},
		{
			ReviewID:     2,
			AverageScore: 8,
			Country:      "US",
			Name:         "Jane Smith",
			Headline:     "Good experience",
			Pros:         "Nice location",
			Cons:         "Could be better",
			Date:         "2024-01-10",
			Language:     "en",
		},
		{
			ReviewID:     3,
			AverageScore: 7,
			Country:      "CA",
			Name:         "Bob Johnson",
			Headline:     "Average stay",
			Pros:         "Decent service",
			Cons:         "Room was small",
			Date:         "2024-01-05",
			Language:     "en",
		},
	}
}

// CreateSampleTranslations creates sample translations for testing
func (td *TestData) CreateSampleTranslations() map[string]*cupid.Property {
	return map[string]*cupid.Property{
		"fr": {
			HotelID:   12345,
			HotelName: "Hôtel de Test Londres",
			Address: cupid.Address{
				City:    "Londres",
				Country: "gb",
			},
		},
		"es": {
			HotelID:   12345,
			HotelName: "Hotel de Prueba Londres",
			Address: cupid.Address{
				City:    "Londres",
				Country: "gb",
			},
		},
		"de": {
			HotelID:   12345,
			HotelName: "Test Hotel London",
			Address: cupid.Address{
				City:    "London",
				Country: "gb",
			},
		},
	}
}

// CreateSampleProperties creates multiple sample properties for testing
func (td *TestData) CreateSampleProperties() []*cupid.Property {
	return []*cupid.Property{
		{
			HotelID:     12345,
			CupidID:     12345,
			HotelName:   "Test Hotel London",
			HotelType:   "Hotels",
			Chain:       "Test Chain",
			Latitude:    51.5074,
			Longitude:   -0.1278,
			Stars:       5,
			Rating:      9.5,
			ReviewCount: 100,
			Address: cupid.Address{
				Address:    "123 Test Street",
				City:       "London",
				State:      "England",
				Country:    "gb",
				PostalCode: "SW1A 1AA",
			},
			MainImageTh: "https://example.com/image1.jpg",
		},
		{
			HotelID:     12346,
			CupidID:     12346,
			HotelName:   "Test Hotel Paris",
			HotelType:   "Hotels",
			Chain:       "Test Chain",
			Latitude:    48.8566,
			Longitude:   2.3522,
			Stars:       4,
			Rating:      8.5,
			ReviewCount: 75,
			Address: cupid.Address{
				Address:    "456 Test Avenue",
				City:       "Paris",
				State:      "Île-de-France",
				Country:    "fr",
				PostalCode: "75001",
			},
			MainImageTh: "https://example.com/image2.jpg",
		},
		{
			HotelID:     12347,
			CupidID:     12347,
			HotelName:   "Test Hotel New York",
			HotelType:   "Hotels",
			Chain:       "Test Chain",
			Latitude:    40.7128,
			Longitude:   -74.0060,
			Stars:       3,
			Rating:      7.5,
			ReviewCount: 50,
			Address: cupid.Address{
				Address:    "789 Test Boulevard",
				City:       "New York",
				State:      "NY",
				Country:    "us",
				PostalCode: "10001",
			},
			MainImageTh: "https://example.com/image3.jpg",
		},
	}
}

// CreateSamplePropertyFilters creates sample property filters for testing
func (td *TestData) CreateSamplePropertyFilters() store.PropertyFilters {
	return store.PropertyFilters{
		City:      "London",
		Country:   "gb",
		MinStars:  3,
		MaxStars:  5,
		MinRating: 7.0,
		MaxRating: 10.0,
		HotelType: "Hotels",
		Chain:     "Test Chain",
	}
}

// CreateEmptyPropertyFilters creates empty property filters for testing
func (td *TestData) CreateEmptyPropertyFilters() store.PropertyFilters {
	return store.PropertyFilters{}
}

// CreateSampleAddress creates a sample address for testing
func (td *TestData) CreateSampleAddress() cupid.Address {
	return cupid.Address{
		Address:    "123 Test Street",
		City:       "London",
		State:      "England",
		Country:    "gb",
		PostalCode: "SW1A 1AA",
	}
}

// CreateSampleFacility creates a sample facility for testing
func (td *TestData) CreateSampleFacility() cupid.Facility {
	return cupid.Facility{
		FacilityID: 1,
		Name:       "WiFi",
	}
}

// CreateSampleFacilities creates multiple sample facilities for testing
func (td *TestData) CreateSampleFacilities() []cupid.Facility {
	return []cupid.Facility{
		{
			FacilityID: 1,
			Name:       "WiFi",
		},
		{
			FacilityID: 2,
			Name:       "Pool",
		},
		{
			FacilityID: 3,
			Name:       "Gym",
		},
	}
}

// CreateSamplePolicy creates a sample policy for testing
func (td *TestData) CreateSamplePolicy() cupid.Policy {
	return cupid.Policy{
		PolicyType:   "checkin",
		Name:         "Check-in",
		Description:  "Check-in from 3:00 PM",
		ChildAllowed: "yes",
		PetsAllowed:  "no",
		Parking:      "free",
	}
}

// CreateSamplePolicies creates multiple sample policies for testing
func (td *TestData) CreateSamplePolicies() []cupid.Policy {
	return []cupid.Policy{
		{
			PolicyType:   "checkin",
			Name:         "Check-in",
			Description:  "Check-in from 3:00 PM",
			ChildAllowed: "yes",
			PetsAllowed:  "no",
			Parking:      "free",
		},
		{
			PolicyType:   "checkout",
			Name:         "Check-out",
			Description:  "Check-out until 11:00 AM",
			ChildAllowed: "yes",
			PetsAllowed:  "no",
			Parking:      "free",
		},
		{
			PolicyType:   "cancellation",
			Name:         "Cancellation",
			Description:  "Free cancellation until 24 hours before check-in",
			ChildAllowed: "yes",
			PetsAllowed:  "no",
			Parking:      "free",
		},
	}
}

// CreateSampleRoom creates a sample room for testing
func (td *TestData) CreateSampleRoom() cupid.Room {
	return cupid.Room{
		ID:             1,
		RoomName:       "Standard Room",
		Description:    "Comfortable standard room",
		RoomSizeSquare: 25,
		RoomSizeUnit:   "m2",
		HotelID:        "12345",
		MaxAdults:      2,
		MaxChildren:    0,
		MaxOccupancy:   2,
		BedRelation:    "1 double bed",
		BedTypes: []cupid.BedType{
			{
				Quantity: 1,
				BedType:  "double",
				BedSize:  "queen",
			},
		},
		RoomAmenities: []cupid.RoomAmenity{
			{
				AmenitiesID: 1,
				Name:        "WiFi",
				Sort:        1,
			},
		},
		Photos: []cupid.Photo{},
		Views:  []cupid.RoomView{},
	}
}

// CreateSampleRooms creates multiple sample rooms for testing
func (td *TestData) CreateSampleRooms() []cupid.Room {
	return []cupid.Room{
		{
			ID:             1,
			RoomName:       "Standard Room",
			Description:    "Comfortable standard room",
			RoomSizeSquare: 25,
			RoomSizeUnit:   "m2",
			HotelID:        "12345",
			MaxAdults:      2,
			MaxChildren:    0,
			MaxOccupancy:   2,
			BedRelation:    "1 double bed",
			BedTypes: []cupid.BedType{
				{
					Quantity: 1,
					BedType:  "double",
					BedSize:  "queen",
				},
			},
			RoomAmenities: []cupid.RoomAmenity{},
			Photos:        []cupid.Photo{},
			Views:         []cupid.RoomView{},
		},
		{
			ID:             2,
			RoomName:       "Deluxe Room",
			Description:    "Spacious deluxe room",
			RoomSizeSquare: 40,
			RoomSizeUnit:   "m2",
			HotelID:        "12345",
			MaxAdults:      4,
			MaxChildren:    2,
			MaxOccupancy:   4,
			BedRelation:    "1 king bed",
			BedTypes: []cupid.BedType{
				{
					Quantity: 1,
					BedType:  "king",
					BedSize:  "king",
				},
			},
			RoomAmenities: []cupid.RoomAmenity{},
			Photos:        []cupid.Photo{},
			Views:         []cupid.RoomView{},
		},
	}
}

// CreateSamplePhoto creates a sample photo for testing
func (td *TestData) CreateSamplePhoto() cupid.Photo {
	return cupid.Photo{
		URL:              "https://example.com/photo1.jpg",
		HDURL:            "https://example.com/photo1_hd.jpg",
		ImageDescription: "Hotel exterior",
		ImageClass1:      "exterior",
		ImageClass2:      "building",
		MainPhoto:        true,
		Score:            0.95,
		ClassID:          1,
		ClassOrder:       1,
	}
}

// CreateSamplePhotos creates multiple sample photos for testing
func (td *TestData) CreateSamplePhotos() []cupid.Photo {
	return []cupid.Photo{
		{
			URL:              "https://example.com/photo1.jpg",
			HDURL:            "https://example.com/photo1_hd.jpg",
			ImageDescription: "Hotel exterior",
			ImageClass1:      "exterior",
			ImageClass2:      "building",
			MainPhoto:        true,
			Score:            0.95,
			ClassID:          1,
			ClassOrder:       1,
		},
		{
			URL:              "https://example.com/photo2.jpg",
			HDURL:            "https://example.com/photo2_hd.jpg",
			ImageDescription: "Lobby",
			ImageClass1:      "interior",
			ImageClass2:      "lobby",
			MainPhoto:        false,
			Score:            0.90,
			ClassID:          2,
			ClassOrder:       2,
		},
		{
			URL:              "https://example.com/photo3.jpg",
			HDURL:            "https://example.com/photo3_hd.jpg",
			ImageDescription: "Room",
			ImageClass1:      "room",
			ImageClass2:      "bedroom",
			MainPhoto:        false,
			Score:            0.85,
			ClassID:          3,
			ClassOrder:       3,
		},
	}
}

// CreateSampleCheckIn creates a sample check-in for testing
func (td *TestData) CreateSampleCheckIn() cupid.CheckIn {
	return cupid.CheckIn{
		CheckInStart:        "15:00",
		CheckInEnd:          "23:00",
		Checkout:            "11:00",
		Instructions:        []cupid.Instruction{},
		SpecialInstructions: "Please contact reception for late check-in",
	}
}

// CreateSamplePropertyWithDetails creates a property with all details for testing
func (td *TestData) CreateSamplePropertyWithDetails() *cupid.Property {
	property := td.CreateSampleProperty()
	property.Facilities = td.CreateSampleFacilities()
	property.Policies = td.CreateSamplePolicies()
	property.Rooms = td.CreateSampleRooms()
	property.Photos = td.CreateSamplePhotos()
	property.CheckIn = td.CreateSampleCheckIn()
	property.Description = "A beautiful test hotel in the heart of London"
	property.MarkdownDescription = "# Test Hotel London\n\nA beautiful test hotel in the heart of London"
	property.ImportantInfo = "Important information about the hotel"
	property.Parking = stringPtr("Free parking available")
	property.GroupRoomMin = intPtr(10)
	property.ChildAllowed = boolPtr(true)
	property.PetsAllowed = boolPtr(false)

	return property
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
