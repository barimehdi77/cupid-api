package cupid

import (
	"time"
)

// Property represents a hotel property from Cupid API
type Property struct {
	HotelID             int64      `json:"hotel_id"`
	CupidID             int64      `json:"cupid_id"`
	MainImageTh         string     `json:"main_image_th"`
	HotelType           string     `json:"hotel_type"`
	HotelTypeID         int        `json:"hotel_type_id"`
	Chain               string     `json:"chain"`
	ChainID             int        `json:"chain_id"`
	Latitude            float64    `json:"latitude"`
	Longitude           float64    `json:"longitude"`
	HotelName           string     `json:"hotel_name"`
	Phone               string     `json:"phone"`
	Fax                 string     `json:"fax"`
	Email               string     `json:"email"`
	Address             Address    `json:"address"`
	Stars               int        `json:"stars"`
	AirportCode         string     `json:"airport_code"`
	Rating              float64    `json:"rating"`
	ReviewCount         int        `json:"review_count"`
	CheckIn             CheckIn    `json:"checkin"`
	Parking             *string    `json:"parking"`
	GroupRoomMin        *int       `json:"group_room_min"`
	ChildAllowed        *bool      `json:"child_allowed"`
	PetsAllowed         *bool      `json:"pets_allowed"`
	Photos              []Photo    `json:"photos"`
	Description         string     `json:"description"`
	MarkdownDescription string     `json:"markdown_description"`
	ImportantInfo       string     `json:"important_info"`
	Facilities          []Facility `json:"facilities"`
	Policies            []Policy   `json:"policies"`
	Rooms               []Room     `json:"rooms"`
	Reviews             *[]Review  `json:"reviews"`
}

// Address represents the hotel address
type Address struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// CheckIn represents check-in information
type CheckIn struct {
	CheckInStart        string        `json:"checkin_start"`
	CheckInEnd          string        `json:"checkin_end"`
	Checkout            string        `json:"checkout"`
	Instructions        []Instruction `json:"instructions"`
	SpecialInstructions string        `json:"special_instructions"`
}

// Instruction represents check-in instructions
type Instruction struct {
	ID          int    `json:"id"`
	Instruction string `json:"instruction"`
}

// Photo represents hotel photos
type Photo struct {
	URL              string  `json:"url"`
	HDURL            string  `json:"hd_url"`
	ImageDescription string  `json:"image_description"`
	ImageClass1      string  `json:"image_class1"`
	ImageClass2      string  `json:"image_class2"`
	MainPhoto        bool    `json:"main_photo"`
	Score            float64 `json:"score"`
	ClassID          int     `json:"class_id"`
	ClassOrder       int     `json:"class_order"`
}

// Facility represents hotel facilities
type Facility struct {
	FacilityID int    `json:"facility_id"`
	Name       string `json:"name"`
}

// Policy represents hotel policies
type Policy struct {
	PolicyType   string `json:"policy_type"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ChildAllowed string `json:"child_allowed"`
	PetsAllowed  string `json:"pets_allowed"`
	Parking      string `json:"parking"`
}

// Room represents hotel room information
type Room struct {
	ID             int64         `json:"id"`
	RoomName       string        `json:"room_name"`
	Description    string        `json:"description"`
	RoomSizeSquare int           `json:"room_size_square"`
	RoomSizeUnit   string        `json:"room_size_unit"`
	HotelID        string        `json:"hotel_id"`
	MaxAdults      int           `json:"max_adults"`
	MaxChildren    int           `json:"max_children"`
	MaxOccupancy   int           `json:"max_occupancy"`
	BedRelation    string        `json:"bed_relation"`
	BedTypes       []BedType     `json:"bed_types"`
	RoomAmenities  []RoomAmenity `json:"room_amenities"`
	Photos         []Photo       `json:"photos"`
	Views          []RoomView    `json:"views"`
}

// BedType represents bed type information
type BedType struct {
	Quantity int    `json:"quantity"`
	BedType  string `json:"bed_type"`
	BedSize  string `json:"bed_size"`
}

// RoomAmenity represents room amenities
type RoomAmenity struct {
	AmenitiesID int    `json:"amenities_id"`
	Name        string `json:"name"`
	Sort        int    `json:"sort"`
}

// RoomView represents room views
type RoomView struct {
	ID   int    `json:"id"`
	View string `json:"view"`
}

// Location represents geographical coordinates (kept for backward compatibility)
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

// ContactInfo represents contact information (kept for backward compatibility)
type ContactInfo struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

// Review represents a property review
type Review struct {
	ReviewID     int64  `json:"review_id"`
	AverageScore int    `json:"average_score"`
	Country      string `json:"country"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Date         string `json:"date"`
	Headline     string `json:"headline"`
	Language     string `json:"language"`
	Pros         string `json:"pros"`
	Cons         string `json:"cons"`
	Source       string `json:"source"`
}

// TranslationResponse represents the translation API response
type TranslationResponse struct {
	Data Property `json:"data"`
}

// Translation represents property translations (kept for backward compatibility)
type Translation struct {
	PropertyID   int64             `json:"property_id"`
	Language     string            `json:"language"`
	Fields       map[string]string `json:"fields"`
	Quality      float64           `json:"quality"`
	TranslatedAt time.Time         `json:"translated_at"`
}

// PropertyData combines all property information
type PropertyData struct {
	Property     Property             `json:"property"`
	Reviews      []Review             `json:"reviews"`
	Translations map[string]*Property `json:"translations"`
}

// PropertyIDs contains all the property IDs from the assignment
var PropertyIDs = []int64{
	1641879, 317597, 1202743, 1037179, 1154868, 1270324, 1305326, 1617655,
	1975211, 2017823, 1503950, 1033299, 378772, 1563003, 1085875, 828917,
	830417, 838887, 1702062, 1144294, 1738870, 898052, 906450, 906467,
	2241195, 1244595, 1277032, 956026, 957111, 152896, 896868, 982911,
	986491, 986622, 988544, 989315, 989544, 990223, 990341, 990370,
	990490, 990609, 990629, 1259611, 991819, 992027, 992851, 993851,
	994085, 994333, 994495, 994903, 995227, 995787, 996977, 1186578,
	999444, 1000017, 1000051, 1198750, 1001100, 1001296, 1001402, 1002200,
	1003142, 1004288, 1006404, 1006602, 1006810, 1006887, 1007101, 1007269,
	1007466, 1011203, 1011644, 1011945, 1012047, 1012140, 1012944, 1023527,
	1013529, 1013584, 1014383, 1015094, 1016591, 1016611, 1017019, 1017039,
	1017044, 1018030, 1018130, 1018251, 1018402, 1018946, 1019473, 1020332,
	1020335, 1020386, 1021856, 1022380,
}
