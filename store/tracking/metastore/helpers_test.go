package metastore_test

import (
	"Muromachi/store/entities"
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"time"
)

// Mock connection with errors (Meta table)
type mockMetaConnectionErrors struct {
}

func (m mockMetaConnectionErrors) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockMetaConnectionErrors) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockMetaErrorRow{}
}

func (m mockMetaConnectionErrors) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, pgx.ErrNoRows
}

type mockMetaErrorRow struct {
}

func (mr mockMetaErrorRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockMetaErrorRow) RawValues() [][]byte {
	return nil
}

func (mr mockMetaErrorRow) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

// Mock connection with successful returned objects (Meta table)
type mockMetaConnection struct {
}

func (m mockMetaConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockMetaConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	t, _ := time.Parse("2006-01-02", "2021-01-19")
	return mockMetaRow{
		Id:               10,
		BundleId:         12,
		Title:            "Im title",
		Price:            "",
		Picture:          "http://picture",
		Screenshots:      []string{"http://picture", "http://picture1", "http://picture"},
		Rating:           "4.6+",
		ReviewCount:      "1002323",
		RatingHistogram:  []string{"1", "2", "3", "4", "5"},
		Description:      "some description of apprepo",
		ShortDescription: "some short description",
		RecentChanges:    "last changes",
		ReleaseDate:      "2020-01-01",
		LastUpdateDate:   "2020-03-03",
		AppSize:          "90MB+",
		Installs:         "1000000+",
		Version:          "v1.3.12",
		AndroidVersion:   "9.0",
		ContentRating:    "18+",
		DeveloperContacts: entities.DeveloperContacts{
			Email:    "email@email.com",
			Contacts: "virginia",
		},
		PrivacyPolicy:  "http://privacypolicy.com",
		Date:           t,
		AppId:          12,
		AppBundle:      "123",
		AppCategory:    "FINANCE",
		AppDeveloperId: "com.develoeper",
		AppDeveloper:   "super developer",
		AppGeo:         "ru_RU",
		AppStartAt:     t.AddDate(-1, 0, 0),
		AppPeriod:      31,
	}
}

func (m mockMetaConnection) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	t, _ := time.Parse("2006-01-02", "2021-01-19")
	for i := 0; i < 3; i++ {
		_ = mockMetaRow{
			Id:               10 + i,
			BundleId:         12,
			Title:            "Im title",
			Price:            "",
			Picture:          "http://picture",
			Screenshots:      []string{"http://picture", "http://picture1", "http://picture"},
			Rating:           "4.6+",
			ReviewCount:      "1002323",
			RatingHistogram:  []string{"1", "2", "3", "4", "5"},
			Description:      "some description of apprepo",
			ShortDescription: "some short description",
			RecentChanges:    "last changes",
			ReleaseDate:      "2020-01-01",
			LastUpdateDate:   "2020-03-03",
			AppSize:          "90MB+",
			Installs:         "1000000+",
			Version:          "v1.3.12",
			AndroidVersion:   "9.0",
			ContentRating:    "18+",
			DeveloperContacts: entities.DeveloperContacts{
				Email:    "email@email.com",
				Contacts: "virginia",
			},
			PrivacyPolicy: "http://privacypolicy.com",
			Date:          t,

			AppId:          12,
			AppBundle:      "123",
			AppCategory:    "FINANCE",
			AppDeveloperId: "com.develoeper",
			AppDeveloper:   "super developer",
			AppGeo:         "ru_RU",
			AppStartAt:     t.AddDate(-1, 0, 0),
			AppPeriod:      31,
		}.Scan(scans...)

		_ = f(mockMetaRow{})

		t = t.AddDate(0, 0, 1)
	}

	return nil, nil
}

type mockMetaRow struct {
	Id                int                     `json:"-"`
	BundleId          int                     `json:"-"`
	Title             string                  `json:"title" db:"title"`
	Price             string                  `json:"price" db:"price"`
	Picture           string                  `json:"picture" db:"picture"`
	Screenshots       []string                `json:"screenshots" db:"screenshots"`
	Rating            string                  `json:"rating" db:"rating"`
	ReviewCount       string                  `json:"reviewCount" db:"review_count"`
	RatingHistogram   []string                `json:"ratingHistogram" db:"rating_histogram"`
	Description       string                  `json:"description" db:"description"`
	ShortDescription  string                  `json:"shortDescription" db:"short_description"`
	RecentChanges     string                  `json:"recentChanges" db:"recent_changes"`
	ReleaseDate       string                  `json:"releaseDate" db:"release_date"`
	LastUpdateDate    string                  `json:"lastUpdateDate" db:"last_update_date"`
	AppSize           string                  `json:"appSize" db:"app_size"`
	Installs          string                  `json:"installs" db:"installs"`
	Version           string                  `json:"version" db:"version"`
	AndroidVersion    string                  `json:"androidVersion" db:"android_version"`
	ContentRating     string                  `json:"contentRating" db:"content_rating"`
	DeveloperContacts entities.DeveloperContacts `json:"developerContacts" db:"developer_contacts"`
	PrivacyPolicy     string                  `json:"privacyPolicy,omitempty"`
	Date              time.Time               `json:"date,omitempty"`
	AppId             int                     `json:"-"`
	AppBundle         string                  `json:"bundle,omitempty"`
	AppCategory       string                  `json:"category,omitempty"`
	AppDeveloperId    string                  `json:"developer_id,omitempty"`
	AppDeveloper      string                  `json:"developer,omitempty"`
	AppGeo            string                  `json:"geo,omitempty"`
	AppStartAt        time.Time               `json:"start_at,omitempty"`
	AppPeriod         uint32                  `json:"period,omitempty"`
}

func (mr mockMetaRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockMetaRow) RawValues() [][]byte {
	return nil
}

func (mr mockMetaRow) Scan(dest ...interface{}) error {
	Id := dest[0].(*int)
	BundleId := dest[1].(*int)
	Title := dest[2].(*string)
	Price := dest[3].(*string)
	Picture := dest[4].(*string)
	Screenshots := dest[5].(*[]string)
	Rating := dest[6].(*string)
	ReviewCount := dest[7].(*string)
	RatingHistogram := dest[8].(*[]string)
	Description := dest[9].(*string)
	ShortDescription := dest[10].(*string)
	RecentChanges := dest[11].(*string)
	ReleaseDate := dest[12].(*string)
	LastUpdateDate := dest[13].(*string)
	AppSize := dest[14].(*string)
	Installs := dest[15].(*string)
	Version := dest[16].(*string)
	AndroidVersion := dest[17].(*string)
	ContentRating := dest[18].(*string)
	DeveloperContacts := dest[19].(*entities.DeveloperContacts)
	PrivacyPolicy := dest[20].(*string)
	Date := dest[21].(*time.Time)
	AppId := dest[22].(*int)
	AppBundle := dest[23].(*string)
	AppCategory := dest[24].(*string)
	AppDeveloperId := dest[25].(*string)
	AppDeveloper := dest[26].(*string)
	AppGeo := dest[27].(*string)
	AppStartAt := dest[28].(*time.Time)
	AppPeriod := dest[29].(*uint32)

	*Id = mr.Id
	*BundleId = mr.BundleId
	*Title = mr.Title
	*Price = mr.Price
	*Picture = mr.Picture
	*Screenshots = mr.Screenshots
	*Rating = mr.Rating
	*ReviewCount = mr.ReviewCount
	*RatingHistogram = mr.RatingHistogram
	*Description = mr.Description
	*ShortDescription = mr.ShortDescription
	*RecentChanges = mr.RecentChanges
	*ReleaseDate = mr.ReleaseDate
	*LastUpdateDate = mr.LastUpdateDate
	*AppSize = mr.AppSize
	*Installs = mr.Installs
	*Version = mr.Version
	*AndroidVersion = mr.AndroidVersion
	*ContentRating = mr.ContentRating
	*DeveloperContacts = mr.DeveloperContacts
	*PrivacyPolicy = mr.PrivacyPolicy
	*Date = mr.Date
	*AppId = mr.AppId
	*AppBundle = mr.AppBundle
	*AppCategory = mr.AppCategory
	*AppDeveloperId = mr.AppDeveloperId
	*AppDeveloper = mr.AppDeveloper
	*AppGeo = mr.AppGeo
	*AppStartAt = mr.AppStartAt
	*AppPeriod = mr.AppPeriod

	return nil
}
