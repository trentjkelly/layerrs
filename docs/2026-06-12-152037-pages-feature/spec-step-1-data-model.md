# Step 1 Spec — Data Model & Database (Pages Feature)

This document defines the exact deliverables for **Step 1** of the Pages feature implementation plan. It covers the database migration, Go entities, file locations, and acceptance criteria.

## Scope

- Add the database schema for Pages, Page tracks, followers, submissions, and editor notes.
- Add Go entity structs under `backend/internals/entities/`.
- Do **not** implement repository methods, service logic, HTTP handlers, or frontend code in this step.

## File locations

| Deliverable | Path |
|-------------|------|
| Goose migration | `backend/migrations/db/00026_create_pages_tables.sql` |
| Page entity | `backend/internals/entities/Page.go` |
| PageTrack entity | `backend/internals/entities/PageTrack.go` |
| PageFollower entity | `backend/internals/entities/PageFollower.go` |
| PageSubmission entity | `backend/internals/entities/PageSubmission.go` |
| PageTrackNote entity | `backend/internals/entities/PageTrackNote.go` |

> Note: The plan references migration `00025_create_pages_tables.sql`, but `00025_delete_track_color.sql` already exists in this repo. Use `00026` instead.

---

## 1. Migration — `backend/migrations/db/00026_create_pages_tables.sql`

Create the following tables, indexes, and constraints. All `created_at`/`updated_at` columns use `TIMESTAMP WITH TIME ZONE` (aliased as `TIMESTAMPTZ`).

### 1.1 `page`

A Page is owned by one artist (the Page Editor). The owner is nullable so deleted accounts become `Anonymous`.

```sql
CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    editor_id INTEGER REFERENCES artist(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (editor_id, name)
);

CREATE INDEX idx_page_editor_id ON page(editor_id);
```

**Rationale from feature:**
- Duplicate page names per editor are prevented.
- Deleted editors result in `editor_id` being `NULL`, displayed as `Anonymous`.

### 1.2 `page_track`

Tracks that have been approved and added to a Page.

```sql
CREATE TABLE page_track (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    recommender_id INTEGER REFERENCES artist(id) ON DELETE SET NULL,
    UNIQUE (page_id, track_id)
);

CREATE INDEX idx_page_track_page_id_added_at ON page_track(page_id, added_at DESC);
CREATE INDEX idx_page_track_track_id ON page_track(track_id);
```

**Rationale from feature:**
- Feed ordering is reverse-chronological by add date.
- `recommender_id` gives credit to the user who submitted/recommended the track.
- Removing a track from a Page removes the Page tag from the track card (no row = no social proof).

### 1.3 `page_follower`

Users following a Page.

```sql
CREATE TABLE page_follower (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    artist_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    followed_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (page_id, artist_id)
);

CREATE INDEX idx_page_follower_page_id ON page_follower(page_id);
CREATE INDEX idx_page_follower_artist_id ON page_follower(artist_id);
```

**Rationale from feature:**
- A Page must have at least one follower to appear on track cards.
- Follower count drives social-proof ranking.

### 1.4 `page_submission`

Tracks submitted/recommended by followers for editor review.

```sql
CREATE TABLE page_submission (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    submitter_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (page_id, track_id)
);

CREATE INDEX idx_page_submission_page_id ON page_submission(page_id);
CREATE INDEX idx_page_submission_track_id ON page_submission(track_id);
CREATE INDEX idx_page_submission_submitter_id ON page_submission(submitter_id);
```

**Rationale from feature:**
- Recommend and submit are the same action.
- A user cannot submit the same track to the same Page twice.
- Only followers can submit.

### 1.5 `page_track_note`

Editor notes attached to an approved Page track. Notes are per-Page-track and only appear on the Page feed.

```sql
CREATE TABLE page_track_note (
    id SERIAL PRIMARY KEY,
    page_track_id INTEGER NOT NULL REFERENCES page_track(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_page_track_note_page_track_id ON page_track_note(page_track_id);
```

**Rationale from feature:**
- Editors can leave notes on submissions they add.
- Notes display on the Page feed but not on the track's artist page.

### 1.6 Full migration file template

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    editor_id INTEGER REFERENCES artist(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (editor_id, name)
);

CREATE INDEX idx_page_editor_id ON page(editor_id);

CREATE TABLE page_track (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    recommender_id INTEGER REFERENCES artist(id) ON DELETE SET NULL,
    UNIQUE (page_id, track_id)
);

CREATE INDEX idx_page_track_page_id_added_at ON page_track(page_id, added_at DESC);
CREATE INDEX idx_page_track_track_id ON page_track(track_id);

CREATE TABLE page_follower (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    artist_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    followed_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (page_id, artist_id)
);

CREATE INDEX idx_page_follower_page_id ON page_follower(page_id);
CREATE INDEX idx_page_follower_artist_id ON page_follower(artist_id);

CREATE TABLE page_submission (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    submitter_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (page_id, track_id)
);

CREATE INDEX idx_page_submission_page_id ON page_submission(page_id);
CREATE INDEX idx_page_submission_track_id ON page_submission(track_id);
CREATE INDEX idx_page_submission_submitter_id ON page_submission(submitter_id);

CREATE TABLE page_track_note (
    id SERIAL PRIMARY KEY,
    page_track_id INTEGER NOT NULL REFERENCES page_track(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_page_track_note_page_track_id ON page_track_note(page_track_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS page_track_note;
DROP TABLE IF EXISTS page_submission;
DROP TABLE IF EXISTS page_follower;
DROP TABLE IF EXISTS page_track;
DROP TABLE IF EXISTS page;

-- +goose StatementEnd
```

---

## 2. Go entities

Add the following files under `backend/internals/entities/`. Match the existing style: `Id`/`CreatedAt`/`UpdatedAt` fields, `json` tags in camelCase, and `time.Time` from the standard library.

### 2.1 `Page.go`

```go
package entities

import "time"

type Page struct {
	Id          int       `json:"id"`
	EditorId    *int      `json:"editorId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Computed fields populated by repository/service layers
	FollowerCount int  `json:"followerCount,omitempty"`
	IsFollowing   bool `json:"isFollowing,omitempty"`
	IsEditor      bool `json:"isEditor,omitempty"`
}

type PageWithFollowerCount struct {
	Page
	EditorName *string `json:"editorName,omitempty"`
}
```

**Field notes:**
- `EditorId` is a pointer so `NULL` in the database maps cleanly to `nil`.
- `FollowerCount`, `IsFollowing`, and `IsEditor` are computed at read time and omitted when zero/default.

### 2.2 `PageTrack.go`

```go
package entities

import "time"

type PageTrack struct {
	Id            int       `json:"id"`
	PageId        int       `json:"pageId"`
	TrackId       int       `json:"trackId"`
	AddedAt       time.Time `json:"addedAt"`
	RecommenderId *int      `json:"recommenderId,omitempty"`

	// Populated when returning a feed item
	Track *Track `json:"track,omitempty"`
}
```

**Field notes:**
- `RecommenderId` is nullable because tracks can be added directly by the editor.
- The embedded `Track` is optional; repositories may join it in when returning feeds.

### 2.3 `PageFollower.go`

```go
package entities

import "time"

type PageFollower struct {
	Id         int       `json:"id"`
	PageId     int       `json:"pageId"`
	ArtistId   int       `json:"artistId"`
	FollowedAt time.Time `json:"followedAt"`
}
```

### 2.4 `PageSubmission.go`

```go
package entities

import "time"

type PageSubmission struct {
	Id          int       `json:"id"`
	PageId      int       `json:"pageId"`
	TrackId     int       `json:"trackId"`
	SubmitterId int       `json:"submitterId"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`

	// Populated when returning submission feed items
	Track     *Track  `json:"track,omitempty"`
	Submitter *Artist `json:"submitter,omitempty"`
}
```

**Field notes:**
- `Track` and `Submitter` are populated by joins so the submission feed can display names without extra round trips.

### 2.5 `PageTrackNote.go`

```go
package entities

import "time"

type PageTrackNote struct {
	Id          int       `json:"id"`
	PageTrackId int       `json:"pageTrackId"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`
}
```

---

## 3. Acceptance criteria

Verify the following before moving to Step 2 (repository layer):

- [ ] Migration file `backend/migrations/db/00026_create_pages_tables.sql` exists and follows the Goose `+goose Up` / `+goose Down` format used by existing migrations.
- [ ] `goose up` applies the migration successfully against a local database.
- [ ] `goose down` rolls back the migration successfully (all five tables are dropped).
- [ ] All tables exist with the correct columns, types, foreign keys, and indexes.
- [ ] The unique constraints are enforced:
  - [ ] `(editor_id, name)` on `page`
  - [ ] `(page_id, track_id)` on `page_track`
  - [ ] `(page_id, artist_id)` on `page_follower`
  - [ ] `(page_id, track_id)` on `page_submission`
- [ ] Deleting an `artist` row sets `page.editor_id` to `NULL` (not cascade).
- [ ] Deleting a `page` row cascades to all child rows in `page_track`, `page_follower`, `page_submission`, and `page_track_note`.
- [ ] All five entity files exist under `backend/internals/entities/`.
- [ ] All entity structs compile and are importable (run `go build ./...` from `backend/`).
- [ ] `Page.EditorId` and `PageTrack.RecommenderId` are pointer types to represent nullable foreign keys.
- [ ] Computed fields (`FollowerCount`, `IsFollowing`, `IsEditor`, embedded `Track`/`Submitter`) use `omitempty` JSON tags or pointers to avoid leaking zero values.

---

## 4. Common pitfalls

1. **Migration numbering:** Do not overwrite or reuse `00025_delete_track_color.sql`. The Pages migration must be `00026_create_pages_tables.sql`.
2. **Repository path:** The plan references `backend/internals/repository/PageRepository.go`, but existing database repositories live at `backend/internals/repository/database/`. Leave that for Step 2.
3. **ON DELETE behavior:** `page.editor_id` must be `SET NULL`, not `CASCADE`, so deleted artists do not delete their Pages.
4. **Nullable IDs:** Use `*int` in Go for `editor_id` and `recommender_id`; do not use zero as a sentinel.
5. **Social-proof index:** `idx_page_track_track_id` is required for efficient lookups when ranking Pages that include a given track.

---

## 5. Verification commands

```bash
# From the backend directory
cd backend

# Apply the migration
goose -dir migrations/db postgres "$DATABASE_URL" up

# Confirm tables exist
psql "$DATABASE_URL" -c "\dt page*"

# Roll back
goose -dir migrations/db postgres "$DATABASE_URL" down

# Build the backend to verify entities compile
go build ./...
```
