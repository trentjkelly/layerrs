# Pages

A **Page** is a curated, followable collection of tracks owned by a single person (the **Page Editor**).

## Core definition

- A Page is owned by exactly one user, the Page Editor.
- Tracks appear in reverse-chronological order: the most recently added track is at the top.
- Other users can **follow** a Page but cannot add tracks to it.
- Other users can **recommend** or **submit** tracks to the Page Editor.
- The Page Editor reviews submissions and decides which tracks are added.
- Being added to a Page is intentionally exclusive and acts as a prestige signal.

## Creating and owning Pages

- Any user can create a Page.
- A user can own as many Pages as they want.
- Page ownership cannot be transferred to another user.
- If a Page Editor deletes their account, the Page continues to exist, but the owner is shown as **Anonymous** (similar to Reddit).

## Page content

- Only **tracks** can be added to a Page.
- There is no limit to the number of tracks a Page can hold.
- The editor can remove tracks from a Page at any time.
- Removing a track from a Page removes that Page's tag from the track card.
- Tracks cannot be reordered; the order is always reverse-chronological.
- Each track on the Page displays when it was added.

## Recommending tracks

- **Recommend** and **submit** are the same action.
- A user must be following a Page in order to recommend a track to it.
- Recommenders get credit for their recommendation on the Page.
- Notifications for accepted or rejected recommendations are a future feature and are not included in this version.

## Editor submission feed

- The Page Editor sees a **Submission Feed** tab alongside the main Page feed.
- The submission feed is only visible to the Page Editor.
- Each submitted track has a button to add it to the Page feed.
- Clicking the add button shows a confirmation pop-up stating the track name and artist name being added.
- The editor can leave text notes on submissions they add.
- Notes are displayed on the Page feed but are not visible when viewing the track elsewhere (for example, on the artist's page).

## Social proof

When a track is added to a well-followed Page, that Page appears on the track card as a tag with the number of followers for the page, for example:

> [Page Name] · 270k

If a track appears on many Pages, the track card surfaces the top Page by follower count and notes how many other Pages include the track:

> [Page Name] · 270k  +254 other pages

Tapping the entire line opens the track page, where every Page that includes the track is ranked by number of followers.

Social proof rules:

- A Page must have at least **1 follower** to appear on the track card.
- If a track is not on any Page with followers, the track card shows nothing.
- If two Pages have the same follower count, they are ranked alphabetically.
- If a track is on more than one Page, the count is shown as `+# other pages`.

## Discovery and following

- Users usually discover Pages through track cards that list their Pages.
- Following a Page adds that Page's tracks to the user's home feed under a **Following** tab.
- The home feed has a Following tab and a main feed, similar to TikTok's Following and For You tabs.
- Pages are always public.

## Notifications and analytics

- The Page Editor does not get notified when someone follows their Page.
- Followers do not get notified when a new track is added; they see it in their Following feed.
- There are no analytics for the Page Editor in this version.

## Relationship to playlists

A Page is **not** a regular playlist:

| Playlist | Page |
|----------|------|
| Often collaborative | Owned by one editor |
| Order is manually set or shuffled | Reverse-chronological by add date |
| Anyone with permission can add | Only the editor can add |
| Casual collection | Curated, exclusive endorsement |

## User roles

- **Page Editor**: owns the Page, reviews submissions, and decides what gets added.
- **Follower**: subscribes to the Page's feed of tracks.
- **Recommender**: submits tracks for the editor to consider.
