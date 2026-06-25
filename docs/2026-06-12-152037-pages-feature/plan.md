# Pages Feature — Implementation Plan

This plan breaks the Pages feature into verifiable sections. Each section lists concrete deliverables and acceptance criteria.

---

## 1. Data model & database

Add Goose migrations and Go entities for Pages and related relationships.

### 1.1 Tables

Create a new migration file `00025_create_pages_tables.sql` with the following tables:

- **`page`**
  - `id` SERIAL PRIMARY KEY
  - `editor_id` INTEGER NOT NULL REFERENCES artist(id) ON DELETE SET NULL
  - `name` VARCHAR(255) NOT NULL
  - `description` TEXT
  - `created_at` TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  - `updated_at` TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  - UNIQUE constraint on `(editor_id, name)` if duplicate page names per editor should be prevented

- **`page_track`**
  - `id` SERIAL PRIMARY KEY
  - `page_id` INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE
  - `track_id` INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE
  - `added_at` TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  - `recommender_id` INTEGER REFERENCES artist(id) ON DELETE SET NULL
  - UNIQUE(page_id, track_id)
  - Index on `(page_id, added_at DESC)` for feed ordering
  - Index on `track_id` for social-proof lookups

- **`page_follower`**
  - `id` SERIAL PRIMARY KEY
  - `page_id` INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE
  - `artist_id` INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE
  - `followed_at` TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  - UNIQUE(page_id, artist_id)

- **`page_submission`**
  - `id` SERIAL PRIMARY KEY
  - `page_id` INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE
  - `track_id` INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE
  - `submitter_id` INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE
  - `note` TEXT
  - `created_at` TIMESTAMP WITH TIME ZONE DEFAULT NOW()
  - UNIQUE(page_id, track_id)

- **`page_track_note`**
  - `id` SERIAL PRIMARY KEY
  - `page_track_id` INTEGER NOT NULL REFERENCES page_track(id) ON DELETE CASCADE
  - `note` TEXT NOT NULL
  - `created_at` TIMESTAMP WITH TIME ZONE DEFAULT NOW()

### 1.2 Go entities

Add new entity files under `backend/internals/entities/`:

- `Page.go`
- `PageTrack.go`
- `PageFollower.go`
- `PageSubmission.go`
- `PageTrackNote.go`

Each entity should include the fields needed for its JSON representation and any computed fields (e.g., follower count).

### 1.3 Acceptance criteria

- [ ] Migration runs successfully with `goose up` and rolls back cleanly with `goose down`.
- [ ] All new tables exist in the database with correct foreign keys and indexes.
- [ ] New entity structs compile and are importable by repository and service layers.

---

## 2. Repository layer

Add repository files under `backend/internals/repository/`:

- `PageRepository.go`

### 2.1 Required repository methods

**Page CRUD**

- `CreatePage(ctx, editorId, name, description) (*Page, error)`
- `GetPageById(ctx, pageId, viewerArtistId int) (*Page, error)` — include `is_editor` and `is_following` flags
- `GetPagesByEditorId(ctx, editorId int) ([]Page, error)`
- `UpdatePage(ctx, pageId, editorId int, name, description string) (*Page, error)`
- `DeletePage(ctx, pageId, editorId int) error`

**Page tracks**

- `AddTrackToPage(ctx, pageId, trackId, recommenderId int) (*PageTrack, error)`
- `RemoveTrackFromPage(ctx, pageId, trackId, editorId int) error`
- `GetPageTracks(ctx, pageId int, limit, offset int) ([]PageTrack, error)` — ordered by `added_at DESC`
- `GetPageTrackById(ctx, pageTrackId int) (*PageTrack, error)`
- `GetPageTrackCount(ctx, pageId int) (int, error)`

**Page submissions**

- `CreateSubmission(ctx, pageId, trackId, submitterId int, note string) (*PageSubmission, error)`
- `GetSubmissionsByPageId(ctx, pageId int, limit, offset int) ([]PageSubmission, error)`
- `DeleteSubmission(ctx, submissionId, editorId int) error`
- `GetSubmissionByPageAndTrack(ctx, pageId, trackId int) (*PageSubmission, error)`

**Page followers**

- `FollowPage(ctx, pageId, artistId int) error`
- `UnfollowPage(ctx, pageId, artistId int) error`
- `GetPageFollowerCount(ctx, pageId int) (int, error)`
- `IsFollowingPage(ctx, pageId, artistId int) (bool, error)`

**Track page social proof**

- `GetPagesForTrack(ctx, trackId int, limit int) ([]PageWithFollowerCount, error)` — only Pages with ≥ 1 follower, ranked by follower count DESC, then name ASC
- `GetPageCountForTrack(ctx, trackId int) (int, error)` — total Pages with ≥ 1 follower

**Notes**

- `CreatePageTrackNote(ctx, pageTrackId int, note string) (*PageTrackNote, error)`
- `GetPageTrackNotes(ctx, pageTrackId int) ([]PageTrackNote, error)`

### 2.2 Acceptance criteria

- [ ] Each repository method has a unit test or manual verification script.
- [ ] `GetPagesForTrack` returns correct ranking (follower count DESC, name ASC for ties).
- [ ] `GetPageTracks` returns tracks in reverse-chronological order.
- [ ] Removing a Page track also removes the ability for that Page to show on the track card.
- [ ] Deleted accounts result in `editor_id` being NULL; repository returns `editor_id` as nullable and frontend displays "Anonymous".

---

## 3. Service layer

Add service files under `backend/internals/service/`:

- `PageService.go`

### 3.1 Required service methods

Map directly to controller needs:

- `CreatePage(ctx, artistId int, req CreatePageRequest) (*entities.Page, error)`
- `GetPage(ctx, pageId, viewerArtistId int) (*entities.Page, error)`
- `UpdatePage(ctx, pageId, artistId int, req UpdatePageRequest) (*entities.Page, error)`
- `DeletePage(ctx, pageId, artistId int) error`
- `GetEditorPages(ctx, editorId int) ([]entities.Page, error)`

- `AddTrackToPage(ctx, pageId, trackId, editorId int) (*entities.PageTrack, error)`
- `RemoveTrackFromPage(ctx, pageId, trackId, editorId int) error`
- `GetPageFeed(ctx, pageId int, pagination Pagination) ([]entities.PageTrack, error)`

- `SubmitTrack(ctx, pageId, trackId, submitterId int, note string) (*entities.PageSubmission, error)`
- `GetSubmissions(ctx, pageId, editorId int, pagination Pagination) ([]entities.PageSubmission, error)`
- `ApproveSubmission(ctx, submissionId, editorId int, note string) (*entities.PageTrack, error)`

- `FollowPage(ctx, pageId, artistId int) error`
- `UnfollowPage(ctx, pageId, artistId int) error`

- `GetTrackPages(ctx, trackId int) (*TrackPagesResponse, error)` — returns top page + count of other pages, plus full list

### 3.2 Business rules

- Only the Page editor can update, delete, add/remove tracks, or view/approve submissions.
- A user must follow a Page to submit a track.
- A user cannot submit the same track to the same Page twice.
- Adding a track from the submission feed creates a `page_track` record and deletes the matching `page_submission`.
- `GetTrackPages` only includes Pages with at least one follower.

### 3.3 Acceptance criteria

- [ ] Service layer enforces all business rules above.
- [ ] `GetTrackPages` returns the correct top page and total count.
- [ ] Editor-only actions return a forbidden error for non-editors.
- [ ] Duplicate submission attempts return a conflict error.

---

## 4. Controller layer & endpoints

Add controller files under `backend/internals/controller/`:

- `PageController.go`

Wire routes in the main router (likely `backend/cmd/api/main.go` or a dedicated router file).

### 4.1 Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/pages` | Create a new Page |
| GET | `/api/pages/:pageId` | Get Page details |
| PATCH | `/api/pages/:pageId` | Update Page name/description |
| DELETE | `/api/pages/:pageId` | Delete Page |
| GET | `/api/artists/:artistId/pages` | Get all Pages owned by an artist |
| POST | `/api/pages/:pageId/follow` | Follow a Page |
| DELETE | `/api/pages/:pageId/follow` | Unfollow a Page |
| POST | `/api/pages/:pageId/tracks` | Add a track to a Page (editor only) |
| DELETE | `/api/pages/:pageId/tracks/:trackId` | Remove a track from a Page (editor only) |
| GET | `/api/pages/:pageId/feed` | Get Page feed (tracks) |
| GET | `/api/pages/:pageId/submissions` | Get submission feed (editor only) |
| POST | `/api/pages/:pageId/submissions` | Submit a track (followers only) |
| POST | `/api/pages/:pageId/submissions/:submissionId/approve` | Approve a submission (editor only) |
| GET | `/api/tracks/:trackId/pages` | Get all Pages for a track, ranked |

### 4.2 Request/response shapes

Define request/response DTOs in `backend/internals/entities/` or a dedicated DTO package.

Example:

```go
type CreatePageRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type TrackPagesResponse struct {
    TopPage      *PageWithFollowerCount   `json:"topPage"`
    OtherCount   int                      `json:"otherCount"`
    Pages        []PageWithFollowerCount  `json:"pages"`
}
```

### 4.3 Acceptance criteria

- [ ] All endpoints return correct HTTP status codes (201 for create, 200 for success, 403 for forbidden, 404 for not found, 409 for conflict).
- [ ] Authentication is enforced via existing JWT middleware.
- [ ] Endpoints are reachable and return expected JSON in local testing.
- [ ] Track card endpoint (`GET /api/tracks/:trackId/pages`) returns only Pages with ≥ 1 follower.

---

## 5. Frontend — Pages

Add SvelteKit routes and components under `frontend/src/`.

### 5.1 New routes

Create the following routes:

- `/pages/new` — Page creation form
- `/pages/[pageId]` — Page detail with tabs
  - `+page.svelte`
  - `+page.server.ts` or load function to fetch Page data
- `/pages/[pageId]/submissions` — Submission feed tab (editor only)
- `/tracks/[trackId]/pages` — Track page listing all Pages

### 5.2 Page detail view

- Display Page name, description, follower count, and editor info (or "Anonymous").
- Show **Follow / Unfollow** button.
- If the viewer is the editor:
  - Show an **Edit** button.
  - Show a **Submissions** tab.
  - Show **Add to feed** buttons on each submission with a confirmation modal showing track name and artist name.
  - Allow the editor to add a note before confirming.
- If the viewer is not the editor:
  - Show a **Recommend track** button if they follow the Page.

### 5.3 Page feed

- List tracks in reverse-chronological order.
- Show when each track was added.
- Show editor notes on the feed only.
- Editor can remove tracks from the feed.

### 5.4 Page creation

- Form with name and description fields.
- Validation: name is required, max length 255.
- On success, redirect to the new Page.

### 5.5 Track page Pages list

- Display all Pages that include the track, ranked by follower count then name.
- Each row shows Page name, follower count, and editor name or "Anonymous".
- Tapping a row navigates to `/pages/[pageId]`.

### 5.6 Acceptance criteria

- [ ] `/pages/new` creates a Page and redirects to it.
- [ ] `/pages/[pageId]` displays Page info and feed correctly.
- [ ] Editor-only tabs and actions are hidden from non-editors.
- [ ] Submission approval modal shows track name and artist name.
- [ ] Adding a note to a submission displays it on the Page feed.
- [ ] Page feed tracks cannot be reordered.
- [ ] Anonymous editor displays as "Anonymous".

---

## 6. Frontend — Track card social proof

Update existing track card components under `frontend/src/components/`.

### 6.1 Track card changes

- If the track is on at least one Page with ≥ 1 follower:
  - Show `[Page Name] · 270k` for the top-ranked Page.
  - If the track is on more than one Page, show `+N other pages` next to it.
- If the track is on zero Pages with followers, show nothing.
- The entire line is tappable and navigates to `/tracks/[trackId]/pages`.

### 6.2 Data fetching

- Add a helper or store to fetch `/api/tracks/[trackId]/pages`.
- Cache or include this data in existing track list endpoints where possible to avoid N+1 requests.

### 6.3 Acceptance criteria

- [ ] Track card shows the top Page by follower count.
- [ ] Track card shows `+N other pages` when there is more than one Page.
- [ ] Tapping the line navigates to the track Pages list.
- [ ] Track cards show nothing when no qualifying Pages exist.
- [ ] Tie-breaking is alphabetical by Page name.

---

## 7. Frontend — Home feed Following tab

Update the home feed route (`frontend/src/routes/+page.svelte` or equivalent).

### 7.1 Feed tabs

- Add a tab switcher with **For You** and **Following** tabs.
- **For You** keeps the existing recommendation feed.
- **Following** shows tracks added to Pages the user follows.

### 7.2 Following feed data

- Fetch `/api/feed/following` or use a query parameter on the existing feed endpoint.
- Return tracks from followed Pages, ordered by `added_at DESC`.
- Optionally group by Page or show Page context on each track card.

### 7.3 Backend endpoint for Following feed

Add a new endpoint:

- `GET /api/feed/following` — returns tracks from all Pages the authenticated user follows, ordered by most recent addition.

Repository method:

- `GetFollowingFeed(ctx, artistId int, limit, offset int) ([]PageTrack, error)`

### 7.4 Acceptance criteria

- [ ] Home feed has two tabs: **For You** and **Following**.
- [ ] Following tab shows tracks from followed Pages.
- [ ] Tracks appear in reverse-chronological order of when they were added to the Page.
- [ ] Unfollowing a Page removes that Page's tracks from the Following feed.

---

## 8. Integration & testing

### 8.1 End-to-end verification

- [ ] Create a Page as user A.
- [ ] User B follows the Page.
- [ ] User B submits a track.
- [ ] User A approves the submission with a note.
- [ ] Track appears on the Page feed with the note.
- [ ] Track card shows the Page tag.
- [ ] User C sees the Page tag on the track card and taps it to view all Pages for the track.
- [ ] User B sees the new track in their Following feed.
- [ ] User A removes the track; the tag disappears from the track card.

### 8.2 Edge cases

- [ ] Page with deleted editor shows "Anonymous" everywhere.
- [ ] Track on multiple Pages with tied follower counts is ranked alphabetically.
- [ ] Submitting a track to a Page the user does not follow is rejected.
- [ ] Duplicate submission to the same Page is rejected.
- [ ] Non-editor attempts to approve a submission are rejected.

---

## 9. Deployment & ops

- [ ] Run `goose up` in production to apply migration `00025_create_pages_tables.sql`.
- [ ] Verify indexes are created and queries use them.
- [ ] Deploy backend and frontend together to avoid broken UI states.

---

## Suggested implementation order

1. Database migrations and entities
2. Repository layer
3. Service layer
4. Controller layer and endpoints
5. Frontend Page creation and detail views
6. Frontend track card social proof
7. Frontend Following tab
8. End-to-end testing and edge-case fixes
