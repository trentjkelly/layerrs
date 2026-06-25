package service

import (
	"context"
	"fmt"

	"github.com/trentjkelly/layerrs/internals/entities"
	databaseRepository "github.com/trentjkelly/layerrs/internals/repository/database"
)

type PageService struct {
	pageRepo  *databaseRepository.PageRepository
	trackRepo *databaseRepository.TrackDatabaseRepository
}

// Constructor for a new PageService
func NewPageService(pageRepo *databaseRepository.PageRepository, trackRepo *databaseRepository.TrackDatabaseRepository) *PageService {
	return &PageService{pageRepo: pageRepo, trackRepo: trackRepo}
}

// --- Page CRUD ---

// Creates a new Page owned by artistId
func (s *PageService) CreatePage(ctx context.Context, artistId int, req entities.CreatePageRequest) (*entities.Page, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("page name is required")
	}

	page, err := s.pageRepo.CreatePage(ctx, artistId, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return page, nil
}

// Gets a Page by id with viewer-specific computed flags
func (s *PageService) GetPage(ctx context.Context, pageId int, viewerArtistId int) (*entities.Page, error) {
	page, err := s.pageRepo.GetPageById(ctx, pageId, viewerArtistId)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	if page == nil {
		return nil, fmt.Errorf("%w", entities.ErrPageNotFound)
	}

	return page, nil
}

// Updates a Page's name and description
func (s *PageService) UpdatePage(ctx context.Context, pageId int, artistId int, req entities.UpdatePageRequest) (*entities.Page, error) {
	page, err := s.requireEditor(ctx, pageId, artistId)
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, fmt.Errorf("page name is required")
	}

	updatedPage, err := s.pageRepo.UpdatePage(ctx, pageId, artistId, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}
	if updatedPage == nil {
		return nil, fmt.Errorf("%w", entities.ErrPageNotFound)
	}

	// Preserve computed viewer flags from the authorization check
	updatedPage.IsEditor = page.IsEditor
	updatedPage.IsFollowing = page.IsFollowing
	updatedPage.FollowerCount = page.FollowerCount

	return updatedPage, nil
}

// Deletes a Page
func (s *PageService) DeletePage(ctx context.Context, pageId int, artistId int) error {
	_, err := s.requireEditor(ctx, pageId, artistId)
	if err != nil {
		return err
	}

	err = s.pageRepo.DeletePage(ctx, pageId, artistId)
	if err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	return nil
}

// Gets all Pages owned by an editor
func (s *PageService) GetEditorPages(ctx context.Context, editorId int) ([]entities.Page, error) {
	pages, err := s.pageRepo.GetPagesByEditorId(ctx, editorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get editor pages: %w", err)
	}

	if pages == nil {
		pages = []entities.Page{}
	}

	return pages, nil
}

// --- Page tracks ---

// Adds a track to a Page directly by the editor
func (s *PageService) AddTrackToPage(ctx context.Context, pageId int, trackId int, editorId int) (*entities.PageTrack, error) {
	_, err := s.requireEditor(ctx, pageId, editorId)
	if err != nil {
		return nil, err
	}

	pageTrack, err := s.pageRepo.AddTrackToPage(ctx, pageId, trackId, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to add track to page: %w", err)
	}

	return pageTrack, nil
}

// Removes a track from a Page
func (s *PageService) RemoveTrackFromPage(ctx context.Context, pageId int, trackId int, editorId int) error {
	_, err := s.requireEditor(ctx, pageId, editorId)
	if err != nil {
		return err
	}

	err = s.pageRepo.RemoveTrackFromPage(ctx, pageId, trackId, editorId)
	if err != nil {
		return fmt.Errorf("failed to remove track from page: %w", err)
	}

	return nil
}

// Gets the feed of tracks for a Page, including editor notes and track info
func (s *PageService) GetPageFeed(ctx context.Context, pageId int, pagination entities.Pagination) ([]entities.PageTrack, error) {
	pageTracks, err := s.pageRepo.GetPageTracks(ctx, pageId, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get page tracks: %w", err)
	}

	if err := s.populatePageTracks(ctx, pageTracks, 0); err != nil {
		return nil, err
	}

	return pageTracks, nil
}

// --- Submissions ---

// Submits a track to a Page. Requires following the Page.
func (s *PageService) SubmitTrack(ctx context.Context, pageId int, trackId int, submitterId int, note string) (*entities.PageSubmission, error) {
	isFollowing, err := s.pageRepo.IsFollowingPage(ctx, pageId, submitterId)
	if err != nil {
		return nil, fmt.Errorf("failed to check follow status: %w", err)
	}
	if !isFollowing {
		return nil, fmt.Errorf("%w", entities.ErrMustFollowPageToSubmit)
	}

	existing, err := s.pageRepo.GetSubmissionByPageAndTrack(ctx, pageId, trackId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing submission: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w", entities.ErrDuplicateSubmission)
	}

	submission, err := s.pageRepo.CreateSubmission(ctx, pageId, trackId, submitterId, note)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	return submission, nil
}

// Gets the submission feed for a Page. Editor only.
func (s *PageService) GetSubmissions(ctx context.Context, pageId int, editorId int, pagination entities.Pagination) ([]entities.PageSubmission, error) {
	_, err := s.requireEditor(ctx, pageId, editorId)
	if err != nil {
		return nil, err
	}

	submissions, err := s.pageRepo.GetSubmissionsByPageId(ctx, pageId, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get submissions: %w", err)
	}

	if submissions == nil {
		submissions = []entities.PageSubmission{}
	}

	return submissions, nil
}

// Approves a submission, adds the track to the Page, optionally creates a note, and removes the submission
func (s *PageService) ApproveSubmission(ctx context.Context, submissionId int, editorId int, note string) (*entities.PageTrack, error) {
	submission, err := s.pageRepo.GetSubmissionById(ctx, submissionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}
	if submission == nil {
		return nil, fmt.Errorf("%w", entities.ErrSubmissionNotFound)
	}

	_, err = s.requireEditor(ctx, submission.PageId, editorId)
	if err != nil {
		return nil, err
	}

	pageTrack, err := s.pageRepo.AddTrackToPage(ctx, submission.PageId, submission.TrackId, submission.SubmitterId)
	if err != nil {
		return nil, fmt.Errorf("failed to add track to page: %w", err)
	}

	if note != "" {
		_, err = s.pageRepo.CreatePageTrackNote(ctx, pageTrack.Id, note)
		if err != nil {
			return nil, fmt.Errorf("failed to create page track note: %w", err)
		}

		notes, err := s.pageRepo.GetPageTrackNotes(ctx, pageTrack.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get page track notes: %w", err)
		}
		pageTrack.Notes = notes
	}

	err = s.pageRepo.DeleteSubmission(ctx, submissionId, editorId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete submission: %w", err)
	}

	if err := s.populatePageTracks(ctx, []entities.PageTrack{*pageTrack}, editorId); err != nil {
		return nil, err
	}

	return pageTrack, nil
}

// --- Followers ---

// Follows a Page
func (s *PageService) FollowPage(ctx context.Context, pageId int, artistId int) error {
	err := s.pageRepo.FollowPage(ctx, pageId, artistId)
	if err != nil {
		return fmt.Errorf("failed to follow page: %w", err)
	}

	return nil
}

// Unfollows a Page
func (s *PageService) UnfollowPage(ctx context.Context, pageId int, artistId int) error {
	err := s.pageRepo.UnfollowPage(ctx, pageId, artistId)
	if err != nil {
		return fmt.Errorf("failed to unfollow page: %w", err)
	}

	return nil
}

// --- Social proof ---

// Gets Pages that include a track, ranked by follower count. Only Pages with >= 1 follower.
func (s *PageService) GetTrackPages(ctx context.Context, trackId int) (*entities.TrackPagesResponse, error) {
	pages, err := s.pageRepo.GetPagesForTrack(ctx, trackId, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages for track: %w", err)
	}

	count, err := s.pageRepo.GetPageCountForTrack(ctx, trackId)
	if err != nil {
		return nil, fmt.Errorf("failed to get page count for track: %w", err)
	}

	response := &entities.TrackPagesResponse{
		Pages:      pages,
		OtherCount: 0,
	}

	if len(pages) > 0 {
		response.TopPage = &pages[0]
		response.OtherCount = count - 1
	}

	return response, nil
}

// --- Following feed ---

// Gets tracks from all Pages the artist follows
func (s *PageService) GetFollowingFeed(ctx context.Context, artistId int, pagination entities.Pagination) ([]entities.PageTrack, error) {
	pageTracks, err := s.pageRepo.GetFollowingFeed(ctx, artistId, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get following feed: %w", err)
	}

	if err := s.populatePageTracks(ctx, pageTracks, artistId); err != nil {
		return nil, err
	}

	return pageTracks, nil
}

// --- Helpers ---

// Populates TrackInfo and editor notes for a slice of PageTracks.
func (s *PageService) populatePageTracks(ctx context.Context, pageTracks []entities.PageTrack, viewerArtistId int) error {
	if len(pageTracks) == 0 {
		return nil
	}

	trackIds := make([]int, 0, len(pageTracks))
	for _, pt := range pageTracks {
		trackIds = append(trackIds, pt.TrackId)
	}

	trackInfos, err := s.trackRepo.ReadTracksByIds(ctx, trackIds, viewerArtistId)
	if err != nil {
		return fmt.Errorf("failed to read track infos: %w", err)
	}

	trackInfoMap := make(map[int]entities.TrackInfo, len(trackInfos))
	for _, info := range trackInfos {
		trackInfoMap[info.Id] = info
	}

	for i := range pageTracks {
		if info, ok := trackInfoMap[pageTracks[i].TrackId]; ok {
			pageTracks[i].Track = &info
		}

		notes, err := s.pageRepo.GetPageTrackNotes(ctx, pageTracks[i].Id)
		if err != nil {
			return fmt.Errorf("failed to get page track notes: %w", err)
		}
		pageTracks[i].Notes = notes
	}

	return nil
}

// Fetches a page and verifies that artistId is the editor. Returns the page on success.
func (s *PageService) requireEditor(ctx context.Context, pageId int, artistId int) (*entities.Page, error) {
	page, err := s.pageRepo.GetPageById(ctx, pageId, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	if page == nil {
		return nil, fmt.Errorf("%w", entities.ErrPageNotFound)
	}
	if !page.IsEditor {
		return nil, fmt.Errorf("%w", entities.ErrForbidden)
	}

	return page, nil
}

// IsEditor checks whether an artist is the editor of a page without returning the full page.
func (s *PageService) IsEditor(ctx context.Context, pageId int, artistId int) (bool, error) {
	page, err := s.pageRepo.GetPageById(ctx, pageId, artistId)
	if err != nil {
		return false, fmt.Errorf("failed to get page: %w", err)
	}
	if page == nil {
		return false, nil
	}
	return page.IsEditor, nil
}

// IsFollowingPage checks whether an artist follows a page.
func (s *PageService) IsFollowingPage(ctx context.Context, pageId int, artistId int) (bool, error) {
	return s.pageRepo.IsFollowingPage(ctx, pageId, artistId)
}

// GetPageFollowerCount returns the follower count for a page.
func (s *PageService) GetPageFollowerCount(ctx context.Context, pageId int) (int, error) {
	return s.pageRepo.GetPageFollowerCount(ctx, pageId)
}
