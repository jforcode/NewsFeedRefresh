package feedSrv

import (
	"database/sql"
	"errors"
	"log"
)

type Main struct {
	db *sql.DB
}

func Init(db *sql.DB) (*Main, error) {
	if db == nil {
		return nil, errors.New("Invalid parameters")
	}

	main := Main{db: db}
	return &main, nil
}

// TODO: bulk insert
// TODO: return number inserted
// TODO: proper logs
func (main *Main) SaveSources(apiSourceName string, sources []Source) error {
	stmt, err := main.db.Prepare("INSERT INTO source (api_source_name, s_id, name, description, url, category, language, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for index, source := range sources {
		log.Printf("Saving source: %+v", source)
		_, err := stmt.Exec(apiSourceName, source.SourceId, source.Name, source.Description, source.Url, source.Category, source.Language, source.Country)
		if err != nil {
			log.Printf("Error while saving source %d: %s", index, err.Error())
		}
		log.Print("Saved")
	}

	return nil
}

func (main *Main) SaveArticles(apiSourceName string, articles []Article) error {
	stmt, err := main.db.Prepare("INSERT INTO article (api_source_name, source_id, source_name, author, title, description, url, url_to_image, published_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for index, article := range articles {
		log.Printf("Saving article: %+v", article)
		_, err := stmt.Exec(apiSourceName, article.SourceId, article.SourceName, article.Author, article.Title, article.Description, article.Url, article.UrlToImage, article.PublishedAt)
		if err != nil {
			log.Printf("Error while saving article %d: %s", index, err.Error())
		}
		log.Print("Saved")
	}

	return nil
}
