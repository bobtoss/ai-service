package document

type Document struct {
	UserID       string `db:"user_id"`
	DocumentID   string `db:"document_id"`
	DocumentName string `db:"document_name"`
}
