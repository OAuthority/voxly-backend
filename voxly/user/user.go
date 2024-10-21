package user

// This struct represents the information about a particular user
// This can be amended at a later stage if more information is needed
type User struct {
	Id string // unique, global identifier for the user
	Username string // global username for the user
	Password string // the users password
	Name string // the users name, if provided
	Email string // the users email
	RegistrationDate string // the users registration date
	Bot bool // is this user a bot?
	Online bool // is this user online?
	Relationship Relationship // the users relationship with the current user
}

// struct for the user relationship
type Relationship struct {
    Type RelationshipType // The type of relationship (e.g., Friend, Blocked)
}


// RelationshipType represents the type of relationship between two users.
type RelationshipType int

// Enum values for RelationshipType
const (
    None RelationshipType = iota // Default value, neither friends or anything else
    Friend // the two users are friends
    Blocked // the user is blocked by the session user
    Outgoing // the session user has sent a request to the user
	Incoming // the session user has a request from the user
    BlockedByUser // the user in the current session is blocked by this user
)
