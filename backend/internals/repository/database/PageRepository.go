package databaseRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type PageRepository struct {
	db *pgxpool.Pool
}

// Constructor for PageRepository
func NewPageRepository(db *pgxpool.Pool) *PageRepository {
	return &PageRepository{db: db}
}

// Closes the database pool connection
func (r *PageRepository) CloseDB() {
	r.db.Close()
}

// --- Page CRUD ---

// Creates a new Page owned by editorId
func (r *PageRepository) CreatePage(ctx context.Context, editorId int, name string, description string) (*entities.Page, error) {
	query := `INSERT INTO page (editor_id, name, description) VALUES ($1, $2, $3) RETURNING id, editor_id, name, description, created_at, updated_at`
	row := r.db.QueryRow(ctx, query, editorId, name, description)

	page, err := scanPageRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return page, nil
}

// Gets a Page by id with computed flags for the viewer
func (r *PageRepository) GetPageById(ctx context.Context, pageId int, viewerArtistId int) (*entities.Page, error) {
	query := `
		SELECT p.id, p.editor_id, p.name, p.description, p.created_at, p.updated_at,
		       COALESCE(fc.follower_count, 0) AS follower_count,
		       CASE WHEN p.editor_id = $2 THEN true ELSE false END AS is_editor,
		       CASE WHEN pf.artist_id IS NOT NULL THEN true ELSE false END AS is_following
		FROM page p
		LEFT JOIN (
			SELECT page_id, COUNT(*) AS follower_count
			FROM page_follower
			GROUP BY page_id
		) fc ON fc.page_id = p.id
		LEFT JOIN page_follower pf ON pf.page_id = p.id AND pf.artist_id = $2
		WHERE p.id = $1
	`
	row := r.db.QueryRow(ctx, query, pageId, viewerArtistId)

	page, err := scanPageRowWithComputed(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get page by id: %w", err)
	}

	return page, nil
}

// Gets all Pages owned by an editor
func (r *PageRepository) GetPagesByEditorId(ctx context.Context, editorId int) ([]entities.Page, error) {
	query := `
		SELECT p.id, p.editor_id, p.name, p.description, p.created_at, p.updated_at,
		       COALESCE(fc.follower_count, 0) AS follower_count
		FROM page p
		LEFT JOIN (
			SELECT page_id, COUNT(*) AS follower_count
			FROM page_follower
			GROUP BY page_id
		) fc ON fc.page_id = p.id
		WHERE p.editor_id = $1
		ORDER BY p.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, editorId)
	if err != nil {
		return nil, fmt.Errorf("failed to query pages by editor: %w", err)
	}
	defer rows.Close()

	var pages []entities.Page
	for rows.Next() {
		page, err := scanPageRowWithFollowerCount(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page row: %w", err)
		}
		pages = append(pages, *page)
	}

	return pages, nil
}

// Updates a Page's name and description if editorId owns it
func (r *PageRepository) UpdatePage(ctx context.Context, pageId int, editorId int, name string, description string) (*entities.Page, error) {
	query := `UPDATE page SET name=$3, description=$4, updated_at=NOW() WHERE id=$1 AND editor_id=$2 RETURNING id, editor_id, name, description, created_at, updated_at`
	row := r.db.QueryRow(ctx, query, pageId, editorId, name, description)

	page, err := scanPageRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	return page, nil
}

// Deletes a Page if editorId owns it
func (r *PageRepository) DeletePage(ctx context.Context, pageId int, editorId int) error {
	query := `DELETE FROM page WHERE id=$1 AND editor_id=$2 RETURNING id`
	row := r.db.QueryRow(ctx, query, pageId, editorId)

	var deletedId int
	err := row.Scan(&deletedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("failed to delete page: %w", err)
	}

	return nil
}

// --- Page tracks ---

// Adds a track to a Page
func (r *PageRepository) AddTrackToPage(ctx context.Context, pageId int, trackId int, recommenderId int) (*entities.PageTrack, error) {
	var recommenderIdPtr *int
	if recommenderId != 0 {
		recommenderIdPtr = &recommenderId
	}

	query := `INSERT INTO page_track (page_id, track_id, recommender_id) VALUES ($1, $2, $3) RETURNING id, page_id, track_id, added_at, recommender_id`
	row := r.db.QueryRow(ctx, query, pageId, trackId, recommenderIdPtr)

	pageTrack, err := scanPageTrackRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to add track to page: %w", err)
	}

	return pageTrack, nil
}

// Removes a track from a Page if editorId owns the Page
func (r *PageRepository) RemoveTrackFromPage(ctx context.Context, pageId int, trackId int, editorId int) error {
	query := `
		DELETE FROM page_track pt
		USING page p
		WHERE pt.page_id=$1 AND pt.track_id=$2 AND p.id=pt.page_id AND p.editor_id=$3
		RETURNING pt.id
	`
	row := r.db.QueryRow(ctx, query, pageId, trackId, editorId)

	var deletedId int
	err := row.Scan(&deletedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("failed to remove track from page: %w", err)
	}

	return nil
}

// Gets tracks in a Page feed ordered by added_at DESC
func (r *PageRepository) GetPageTracks(ctx context.Context, pageId int, limit int, offset int) ([]entities.PageTrack, error) {
	query := `
		SELECT id, page_id, track_id, added_at, recommender_id
		FROM page_track
		WHERE page_id=$1
		ORDER BY added_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, pageId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query page tracks: %w", err)
	}
	defer rows.Close()

	var pageTracks []entities.PageTrack
	for rows.Next() {
		pageTrack, err := scanPageTrackRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page track row: %w", err)
		}
		pageTracks = append(pageTracks, *pageTrack)
	}

	return pageTracks, nil
}

// Gets a single PageTrack by id
func (r *PageRepository) GetPageTrackById(ctx context.Context, pageTrackId int) (*entities.PageTrack, error) {
	query := `SELECT id, page_id, track_id, added_at, recommender_id FROM page_track WHERE id=$1`
	row := r.db.QueryRow(ctx, query, pageTrackId)

	pageTrack, err := scanPageTrackRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get page track by id: %w", err)
	}

	return pageTrack, nil
}

// Gets the number of tracks in a Page
func (r *PageRepository) GetPageTrackCount(ctx context.Context, pageId int) (int, error) {
	query := `SELECT COUNT(*) FROM page_track WHERE page_id=$1`

	var count int
	err := r.db.QueryRow(ctx, query, pageId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get page track count: %w", err)
	}

	return count, nil
}

// --- Page submissions ---

// Creates a new track submission for a Page
func (r *PageRepository) CreateSubmission(ctx context.Context, pageId int, trackId int, submitterId int, note string) (*entities.PageSubmission, error) {
	query := `
		INSERT INTO page_submission (page_id, track_id, submitter_id, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, page_id, track_id, submitter_id, note, created_at
	`
	row := r.db.QueryRow(ctx, query, pageId, trackId, submitterId, note)

	submission, err := scanPageSubmissionRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	return submission, nil
}

// Gets submissions for a Page ordered by created_at DESC
func (r *PageRepository) GetSubmissionsByPageId(ctx context.Context, pageId int, limit int, offset int) ([]entities.PageSubmission, error) {
	query := `
		SELECT id, page_id, track_id, submitter_id, note, created_at
		FROM page_submission
		WHERE page_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, pageId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []entities.PageSubmission
	for rows.Next() {
		submission, err := scanPageSubmissionRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission row: %w", err)
		}
		submissions = append(submissions, *submission)
	}

	return submissions, nil
}

// Deletes a submission if editorId owns the Page
func (r *PageRepository) DeleteSubmission(ctx context.Context, submissionId int, editorId int) error {
	query := `
		DELETE FROM page_submission ps
		USING page p
		WHERE ps.id=$1 AND ps.page_id=p.id AND p.editor_id=$2
		RETURNING ps.id
	`
	row := r.db.QueryRow(ctx, query, submissionId, editorId)

	var deletedId int
	err := row.Scan(&deletedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("failed to delete submission: %w", err)
	}

	return nil
}

// Gets a submission by Page and track id
func (r *PageRepository) GetSubmissionByPageAndTrack(ctx context.Context, pageId int, trackId int) (*entities.PageSubmission, error) {
	query := `
		SELECT id, page_id, track_id, submitter_id, note, created_at
		FROM page_submission
		WHERE page_id=$1 AND track_id=$2
	`
	row := r.db.QueryRow(ctx, query, pageId, trackId)

	submission, err := scanPageSubmissionRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get submission by page and track: %w", err)
	}

	return submission, nil
}

// Gets a submission by its id
func (r *PageRepository) GetSubmissionById(ctx context.Context, submissionId int) (*entities.PageSubmission, error) {
	query := `
		SELECT id, page_id, track_id, submitter_id, note, created_at
		FROM page_submission
		WHERE id=$1
	`
	row := r.db.QueryRow(ctx, query, submissionId)

	submission, err := scanPageSubmissionRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get submission by id: %w", err)
	}

	return submission, nil
}

// --- Page followers ---

// Adds a follower to a Page
func (r *PageRepository) FollowPage(ctx context.Context, pageId int, artistId int) error {
	query := `INSERT INTO page_follower (page_id, artist_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, pageId, artistId)
	if err != nil {
		return fmt.Errorf("failed to follow page: %w", err)
	}

	return nil
}

// Removes a follower from a Page
func (r *PageRepository) UnfollowPage(ctx context.Context, pageId int, artistId int) error {
	query := `DELETE FROM page_follower WHERE page_id=$1 AND artist_id=$2`
	_, err := r.db.Exec(ctx, query, pageId, artistId)
	if err != nil {
		return fmt.Errorf("failed to unfollow page: %w", err)
	}

	return nil
}

// Gets the follower count for a Page
func (r *PageRepository) GetPageFollowerCount(ctx context.Context, pageId int) (int, error) {
	query := `SELECT COUNT(*) FROM page_follower WHERE page_id=$1`

	var count int
	err := r.db.QueryRow(ctx, query, pageId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get page follower count: %w", err)
	}

	return count, nil
}

// Checks if an artist is following a Page
func (r *PageRepository) IsFollowingPage(ctx context.Context, pageId int, artistId int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM page_follower WHERE page_id=$1 AND artist_id=$2)`

	var isFollowing bool
	err := r.db.QueryRow(ctx, query, pageId, artistId).Scan(&isFollowing)
	if err != nil {
		return false, fmt.Errorf("failed to check page follow status: %w", err)
	}

	return isFollowing, nil
}

// --- Track page social proof ---

// Gets Pages that include a track, only those with >= 1 follower, ranked by followers DESC then name ASC
func (r *PageRepository) GetPagesForTrack(ctx context.Context, trackId int, limit int) ([]entities.PageWithFollowerCount, error) {
	query := `
		SELECT p.id, p.editor_id, p.name, p.description, p.created_at, p.updated_at,
		       COUNT(pf.artist_id) AS follower_count,
		       a.username
		FROM page p
		JOIN page_track pt ON pt.page_id = p.id
		LEFT JOIN page_follower pf ON pf.page_id = p.id
		LEFT JOIN artist a ON a.id = p.editor_id
		WHERE pt.track_id = $1
		GROUP BY p.id, p.editor_id, p.name, p.description, p.created_at, p.updated_at, a.username
		HAVING COUNT(pf.artist_id) >= 1
		ORDER BY follower_count DESC, p.name ASC
	`
	if limit > 0 {
		query += ` LIMIT $2`
	}
	var rows pgx.Rows
	var err error
	if limit > 0 {
		rows, err = r.db.Query(ctx, query, trackId, limit)
	} else {
		rows, err = r.db.Query(ctx, query, trackId)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query pages for track: %w", err)
	}
	defer rows.Close()

	var pages []entities.PageWithFollowerCount
	for rows.Next() {
		page, err := scanPageWithFollowerCountRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page with follower count row: %w", err)
		}
		pages = append(pages, *page)
	}

	return pages, nil
}

// Gets the total number of Pages with >= 1 follower that include a track
func (r *PageRepository) GetPageCountForTrack(ctx context.Context, trackId int) (int, error) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT p.id
			FROM page p
			JOIN page_track pt ON pt.page_id = p.id
			LEFT JOIN page_follower pf ON pf.page_id = p.id
			WHERE pt.track_id = $1
			GROUP BY p.id
			HAVING COUNT(pf.artist_id) >= 1
		) sub
	`

	var count int
	err := r.db.QueryRow(ctx, query, trackId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get page count for track: %w", err)
	}

	return count, nil
}

// --- Notes ---

// Creates a note attached to a PageTrack
func (r *PageRepository) CreatePageTrackNote(ctx context.Context, pageTrackId int, note string) (*entities.PageTrackNote, error) {
	query := `
		INSERT INTO page_track_note (page_track_id, note)
		VALUES ($1, $2)
		RETURNING id, page_track_id, note, created_at
	`
	row := r.db.QueryRow(ctx, query, pageTrackId, note)

	pageTrackNote, err := scanPageTrackNoteRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to create page track note: %w", err)
	}

	return pageTrackNote, nil
}

// Gets notes for a PageTrack
func (r *PageRepository) GetPageTrackNotes(ctx context.Context, pageTrackId int) ([]entities.PageTrackNote, error) {
	query := `
		SELECT id, page_track_id, note, created_at
		FROM page_track_note
		WHERE page_track_id=$1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, pageTrackId)
	if err != nil {
		return nil, fmt.Errorf("failed to query page track notes: %w", err)
	}
	defer rows.Close()

	var notes []entities.PageTrackNote
	for rows.Next() {
		note, err := scanPageTrackNoteRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page track note row: %w", err)
		}
		notes = append(notes, *note)
	}

	return notes, nil
}

// --- Following feed ---

// Gets tracks from all Pages an artist follows, ordered by most recently added
func (r *PageRepository) GetFollowingFeed(ctx context.Context, artistId int, limit int, offset int) ([]entities.PageTrack, error) {
	query := `
		SELECT pt.id, pt.page_id, pt.track_id, pt.added_at, pt.recommender_id
		FROM page_track pt
		JOIN page_follower pf ON pf.page_id = pt.page_id
		WHERE pf.artist_id = $1
		ORDER BY pt.added_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, artistId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query following feed: %w", err)
	}
	defer rows.Close()

	var pageTracks []entities.PageTrack
	for rows.Next() {
		pageTrack, err := scanPageTrackRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan following feed row: %w", err)
		}
		pageTracks = append(pageTracks, *pageTrack)
	}

	return pageTracks, nil
}

// --- Scan helpers ---

type pageScanner interface {
	Scan(dest ...any) error
}

func scanPageRow(scanner pageScanner) (*entities.Page, error) {
	var page entities.Page
	var editorId sql.NullInt32
	var description sql.NullString

	err := scanner.Scan(
		&page.Id,
		&editorId,
		&page.Name,
		&description,
		&page.CreatedAt,
		&page.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if editorId.Valid {
		id := int(editorId.Int32)
		page.EditorId = &id
	}
	if description.Valid {
		page.Description = description.String
	}

	return &page, nil
}

func scanPageRowWithComputed(scanner pageScanner) (*entities.Page, error) {
	var page entities.Page
	var editorId sql.NullInt32
	var description sql.NullString

	err := scanner.Scan(
		&page.Id,
		&editorId,
		&page.Name,
		&description,
		&page.CreatedAt,
		&page.UpdatedAt,
		&page.FollowerCount,
		&page.IsEditor,
		&page.IsFollowing,
	)
	if err != nil {
		return nil, err
	}

	if editorId.Valid {
		id := int(editorId.Int32)
		page.EditorId = &id
	}
	if description.Valid {
		page.Description = description.String
	}

	return &page, nil
}

func scanPageRowWithFollowerCount(scanner pageScanner) (*entities.Page, error) {
	var page entities.Page
	var editorId sql.NullInt32
	var description sql.NullString

	err := scanner.Scan(
		&page.Id,
		&editorId,
		&page.Name,
		&description,
		&page.CreatedAt,
		&page.UpdatedAt,
		&page.FollowerCount,
	)
	if err != nil {
		return nil, err
	}

	if editorId.Valid {
		id := int(editorId.Int32)
		page.EditorId = &id
	}
	if description.Valid {
		page.Description = description.String
	}

	return &page, nil
}

func scanPageTrackRow(scanner pageScanner) (*entities.PageTrack, error) {
	var pageTrack entities.PageTrack
	var recommenderId sql.NullInt32

	err := scanner.Scan(
		&pageTrack.Id,
		&pageTrack.PageId,
		&pageTrack.TrackId,
		&pageTrack.AddedAt,
		&recommenderId,
	)
	if err != nil {
		return nil, err
	}

	if recommenderId.Valid {
		id := int(recommenderId.Int32)
		pageTrack.RecommenderId = &id
	}

	return &pageTrack, nil
}

func scanPageSubmissionRow(scanner pageScanner) (*entities.PageSubmission, error) {
	var submission entities.PageSubmission
	var note sql.NullString

	err := scanner.Scan(
		&submission.Id,
		&submission.PageId,
		&submission.TrackId,
		&submission.SubmitterId,
		&note,
		&submission.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if note.Valid {
		submission.Note = note.String
	}

	return &submission, nil
}

func scanPageTrackNoteRow(scanner pageScanner) (*entities.PageTrackNote, error) {
	var note entities.PageTrackNote

	err := scanner.Scan(
		&note.Id,
		&note.PageTrackId,
		&note.Note,
		&note.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func scanPageWithFollowerCountRow(scanner pageScanner) (*entities.PageWithFollowerCount, error) {
	var pwc entities.PageWithFollowerCount
	var editorId sql.NullInt32
	var description sql.NullString
	var editorName sql.NullString

	err := scanner.Scan(
		&pwc.Id,
		&editorId,
		&pwc.Name,
		&description,
		&pwc.CreatedAt,
		&pwc.UpdatedAt,
		&pwc.FollowerCount,
		&editorName,
	)
	if err != nil {
		return nil, err
	}

	if editorId.Valid {
		id := int(editorId.Int32)
		pwc.EditorId = &id
	}
	if description.Valid {
		pwc.Description = description.String
	}
	if editorName.Valid {
		pwc.EditorName = &editorName.String
	}

	return &pwc, nil
}
