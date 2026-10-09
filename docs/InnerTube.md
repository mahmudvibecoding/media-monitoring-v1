# InnerTube API reference

**YouTube WEB client · Public information retrieval · Updated 2026-10-08**

Reference for channel discovery, uploads, video metadata, views, likes, and
comments. Search, transcripts, Shorts, and other reading endpoints are included
as supplementary capabilities. InnerTube is YouTube's internal API; it has no
stable public schema. The JSON paths below describe the WEB response shapes.

## Contents

- [Endpoint index](#endpoint-index)
- [HTTP requests and client context](#http-requests-and-client-context)
- [Identifiers](#identifiers)
- [Resolve a URL](#resolve-a-url)
- [Read channel information and tabs](#read-channel-information-and-tabs)
- [Read an uploads playlist](#read-an-uploads-playlist)
- [Read video metadata](#read-video-metadata)
- [Identify Shorts and live content](#identify-shorts-and-live-content)
- [Read views and likes](#read-views-and-likes)
- [Read top-level comments](#read-top-level-comments)
- [Search](#search)
- [Read a transcript](#read-a-transcript)
- [Shorts viewer endpoints](#shorts-viewer-endpoints)
- [Live chat](#live-chat)
- [Other reading endpoints](#other-reading-endpoints)
- [Pagination and response parsing](#pagination-and-response-parsing)
- [Errors and data semantics](#errors-and-data-semantics)

## Endpoint index

All paths are relative to `https://www.youtube.com/youtubei/v1` unless a different
YouTube client host is required. The operations below use **POST**.

| Information | Endpoint | Main request fields |
|---|---|---|
| Internal destination for a channel, video, or playlist URL | `/navigation/resolve_url` | `url` |
| Channel metadata, channel tabs, playlist entries | `/browse` | `browseId`, optional `params`; or `continuation` |
| Video title, description, dates, duration, and content flags | `/player` | `videoId` |
| Video view and like count updates | `/updated_metadata` | `videoId`, optional returned `continuation` |
| Watch-page information, comment entry point, comment pages | `/next` | `videoId`; or comment `continuation` |
| Search results | `/search` | `query`, optional filter `params`; or `continuation` |
| Available transcript segments | `/get_transcript` | Returned transcript `params` |
| One video's Shorts-viewer information | `/reel/reel_item_watch` | `playerRequest`, returned `params` |
| Shorts recommendation sequence | `/reel/reel_watch_sequence` | Returned `sequenceParams` |
| Live chat events | `/live_chat/get_live_chat` | Chat `continuation` |
| Archived live chat events | `/live_chat/get_live_chat_replay` | Replay `continuation` |

Additional interface, Music, and authenticated-account endpoints are listed in
[Other reading endpoints](#other-reading-endpoints).

## HTTP requests and client context

### URL, method, and headers

```http
POST https://www.youtube.com/youtubei/v1/updated_metadata?prettyPrint=false
Content-Type: application/json
Origin: https://www.youtube.com
X-Youtube-Client-Name: 1
X-Youtube-Client-Version: 2.20261007.01.00
```

`prettyPrint=false` requests compact JSON. POST is also used for reading data;
it does not inherently mean creating or modifying content.

Use a matching client name and version in the headers and JSON body. For the WEB
client, the header uses numeric identifier `1`, while the body uses `"WEB"`.

### Common body

```json
{
  "context": {
    "client": {
      "clientName": "WEB",
      "clientVersion": "2.20261007.01.00",
      "hl": "en",
      "gl": "US"
    }
  },
  "videoId": "XuhGo5OY3aQ"
}
```

| Field | Meaning |
|---|---|
| `context.client.clientName` | Client type; use `WEB` for these response shapes |
| `context.client.clientVersion` | Matching YouTube client version |
| `context.client.hl` | Preferred interface language; affects labels and display text |
| `context.client.gl` | Requested country context; does not change the connection's IP location |
| `context.client.visitorData` | Optional visitor/session data when supplied by YouTube |

The shown version is a concrete example, not a permanent API version. Obtain a
current value from a WEB request in the browser's Network panel or from the page's
`INNERTUBE_CLIENT_VERSION` configuration. Keep associated session context consistent
when following continuations. If using `X-Goog-Visitor-Id`, keep it consistent with
`context.client.visitorData`. [Client configuration](https://github.com/LuanRT/YouTube.js/blob/main/src/core/Session.ts)

These public WEB examples do not include account credentials or an API key.
Anonymous access can still be restricted, and account endpoints require an
authenticated session. Client identification does not authenticate an account.

**Every request fragment below must be merged with the common `context`.**
Uppercase values such as `VIDEO_ID` or `RETURNED_TOKEN` are placeholders.

### Complete request example

```bash
curl --silent --show-error --fail-with-body --compressed --max-time 30 \
  'https://www.youtube.com/youtubei/v1/updated_metadata?prettyPrint=false' \
  -H 'Content-Type: application/json' \
  -H 'Origin: https://www.youtube.com' \
  -H 'X-Youtube-Client-Name: 1' \
  -H 'X-Youtube-Client-Version: 2.20261007.01.00' \
  --data-binary @- <<'JSON'
{
  "context": {
    "client": {
      "clientName": "WEB",
      "clientVersion": "2.20261007.01.00",
      "hl": "en",
      "gl": "US"
    }
  },
  "videoId": "XuhGo5OY3aQ"
}
JSON
```

The endpoint returns JSON. A browser URL such as `/watch?v=...` or
`/playlist?list=...` normally returns an HTML page instead.

## Identifiers

| Identifier | Example | Meaning |
|---|---|---|
| Channel ID | `UCDPnSqLmYqHadD6VkrRiarg` | Stable channel identifier |
| Handle | `@RizaNova` | Human-readable channel address; can change |
| Video ID | `XuhGo5OY3aQ` | Identifies a regular video, Short, or livestream |
| Uploads playlist ID | `UUDPnSqLmYqHadD6VkrRiarg` | Playlist associated with channel uploads |
| Playlist browse ID | `VLUUDPnSqLmYqHadD6VkrRiarg` | Playlist destination accepted by `/browse` |
| Continuation | Opaque string returned in a response | Cursor for a specific collection or update stream |
| `params` | Opaque string returned in a navigation command | Selects a tab, filter, panel, or other view |

For standard channel uploads, the channel's `UC` prefix is commonly replaced with
`UU`; prepend `VL` to the resulting playlist ID for `/browse`. Validate the returned
playlist rather than assuming every derived ID is available. Other playlist
prefixes serve other purposes. Do not replace a video ID's characters to classify
its content type.

Tokens and encoded parameters are not ordinary identifiers. Preserve them exactly,
including `%3D` and similar sequences; do not decode, double-encode, or reuse them
for another channel, sort order, endpoint, or session.

## Resolve a URL

**Endpoint:** `/navigation/resolve_url`

```json
{
  "url": "https://www.youtube.com/@RizaNova"
}
```

For a channel, inspect:

```text
endpoint.browseEndpoint.browseId
endpoint.browseEndpoint.params
endpoint.commandMetadata.webCommandMetadata.apiUrl
```

A video URL may resolve to a watch destination instead. Inspect the returned
endpoint type before reading its identifier. A successful HTTP response without
an appropriate destination is not a successful channel resolution.
[URL resolution](https://github.com/LuanRT/YouTube.js/blob/main/src/Innertube.ts)

## Read channel information and tabs

**Endpoint:** `/browse`

```json
{
  "browseId": "UCDPnSqLmYqHadD6VkrRiarg"
}
```

Channel metadata is commonly under `metadata.channelMetadataRenderer`:

| Relative field | Information |
|---|---|
| `externalId` | Channel ID |
| `title` | Channel name |
| `description` | Channel description |
| `avatar.thumbnails[]` | Avatar URLs and dimensions |
| `channelUrl` | Channel URL |
| `vanityChannelUrl` | Handle-based URL when available |
| `keywords` | Channel keywords when exposed |

The page header can also contain displayed subscriber and video counts. Header
layouts vary; abbreviated text such as `15M subscribers` is not an exact integer.

Read available tabs from:

```text
contents.twoColumnBrowseResultsRenderer.tabs[].tabRenderer
```

Each tab can expose `title`, `selected`, and
`endpoint.browseEndpoint.{browseId, params}`. To request Videos, Shorts, Live,
Playlists, or Posts, use the endpoint supplied for that tab:

```json
{
  "browseId": "CHANNEL_ID",
  "params": "PARAMS_FROM_THE_SELECTED_TAB"
}
```

Tab labels are localized, and tabs may be absent. Prefer the endpoint's destination
URL, when provided, over relying solely on an English title or a fixed array index.
[Channel handling](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/Channel.ts),
[browse request fields](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/endpoints/BrowseEndpoint.ts)

## Read an uploads playlist

**Endpoint:** `/browse`

```json
{
  "browseId": "VLUUDPnSqLmYqHadD6VkrRiarg"
}
```

### Video entries

Within the selected playlist tab, video cards can appear at:

```text
contents.twoColumnBrowseResultsRenderer.tabs[].tabRenderer.content
  .sectionListRenderer.contents[].itemSectionRenderer.contents[]
  .lockupViewModel
```

Inspect the selected tab's content rather than every card anywhere in the response.

| Field relative to `lockupViewModel` | Information |
|---|---|
| `contentId` | Video ID |
| `contentType` | `LOCKUP_CONTENT_TYPE_VIDEO` for a video card |
| `metadata.lockupMetadataViewModel.title.content` | Title |
| `contentImage.thumbnailViewModel.image.sources[]` | Thumbnail URLs and dimensions |
| `metadata.lockupMetadataViewModel.metadata.contentMetadataViewModel.metadataRows[]` | Display information such as uploader, views, and relative age |

A response can contain repeated video IDs inside queue, sharing, and navigation
controls. Read content entries and deduplicate by `contentId`; recursively collecting
every `videoId` does not produce a reliable playlist.

Uploads can include regular videos and Shorts. `LOCKUP_CONTENT_TYPE_VIDEO` does
not distinguish them. Playlist display metadata is also not a substitute for a
full video description or a precise publication timestamp.

### Next page

A sibling `continuationItemViewModel` can provide:

```text
continuationItemViewModel.continuationCommand.innertubeCommand
  .continuationCommand.token
```

Submit the token to the endpoint indicated by that command:

```json
{
  "continuation": "UPLOADS_CONTINUATION_TOKEN"
}
```

The next page can put its items under:

```text
onResponseReceivedActions[].appendContinuationItemsAction.continuationItems[]
```

Read the cards and any following continuation from that array. Page sizes are not
a contract. An absent continuation ends the returned collection; it does not prove
that private, removed, or otherwise unavailable uploads do not exist.

## Read video metadata

**Endpoint:** `/player`

```json
{
  "videoId": "XuhGo5OY3aQ"
}
```

### Main fields

| Response path | Information and type |
|---|---|
| `videoDetails.videoId` | Video ID, string |
| `videoDetails.channelId` | Uploader's channel ID, string |
| `videoDetails.title` | Video title, string |
| `videoDetails.shortDescription` | Description text, string; the name does not imply a short excerpt |
| `videoDetails.author` | Uploader's displayed name |
| `videoDetails.lengthSeconds` | Duration, commonly a decimal string |
| `videoDetails.viewCount` | Video views, commonly a decimal string |
| `videoDetails.keywords[]` | Video keywords, when present |
| `videoDetails.thumbnail.thumbnails[]` | Thumbnail URL, width, and height |
| `videoDetails.isLive` | Whether the response identifies an active live broadcast |
| `videoDetails.isUpcoming` | Upcoming-content indicator |
| `videoDetails.isLiveContent` | Live-related content, including content that is no longer live |

Missing fields remain unknown. Numeric strings should be validated before conversion.
[Video detail fields](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/misc/VideoDetails.ts)

### Publication information and additional metadata

Use `microformat.playerMicroformatRenderer` as the prefix for these fields:

| Relative field | Information |
|---|---|
| `publishDate` | Publication date or timestamp |
| `uploadDate` | Upload date or timestamp |
| `category` | Video category |
| `externalChannelId` | Channel ID |
| `ownerChannelName` | Displayed channel name |
| `ownerProfileUrl` | Channel profile URL |
| `canonicalUrl` | Canonical video URL when provided |
| `isShortsEligible` | Shorts eligibility signal when provided |
| `likeCount` | Numeric-string likes when exposed in this response |
| `viewCount` | Additional view-count representation |
| `liveBroadcastDetails` | Live start/end information when available |

Upload and publication dates are separate concepts. Preserve supplied timezone
offsets; date-only values do not specify an exact moment. Prefer
`videoDetails.lengthSeconds` for the main duration rather than combining slightly
different representations. [Microformat structure](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/PlayerMicroformat.ts)

### Playback information

`playabilityStatus` describes playback eligibility. `streamingData` can describe
media formats and URLs, while `captions.playerCaptionsTracklistRenderer.captionTracks[]`
can describe available caption tracks. These sections are optional; stream bytes
and transcript text are not embedded in the video details.

**Playback failure and metadata failure are different.** A response marked
`UNPLAYABLE` can still contain usable metadata. A response with `LOGIN_REQUIRED`
and no `videoDetails` cannot supply the requested video information. Validate the
fields needed for the operation instead of treating HTTP 200 or playback status
alone as the result. [Player-response handling](https://github.com/LuanRT/YouTube.js/blob/main/src/core/mixins/MediaInfo.ts)

## Identify Shorts and live content

A video ID has the same format across regular videos, Shorts, and livestreams.

| Signal | Interpretation |
|---|---|
| `isShortsEligible: true` plus a matching `/shorts/VIDEO_ID` canonical URL | Strong, mutually supporting Shorts signals |
| `isShortsEligible: false` with a watch URL | Supports regular-video classification |
| A video entry in the channel's Shorts tab | Evidence that YouTube presents the item as a Short |
| `LOCKUP_CONTENT_TYPE_VIDEO` | Generic video content; insufficient to identify a Short |
| Short duration or portrait dimensions | Insufficient alone |
| `isLiveContent: true` | Live-related content; does not mean currently live |

Preserve the raw signals. If they are missing or inconsistent, classification is
unknown. Do not infer classification from recommendations for other videos.

On a channel's Shorts tab, cards can use `shortsLockupViewModel`, with a video ID at:

```text
shortsLockupViewModel.onTap.innertubeCommand.reelWatchEndpoint.videoId
```

Its `entityId` is a UI entity identifier, not necessarily the plain video ID.

## Read views and likes

### Endpoint choice

| Endpoint | Views | Likes | Response form |
|---|---|---|---|
| `/updated_metadata` | Viewership update action | Numeric fields in a like-count entity when available | Counter updates and entity mutations |
| `/player` | `videoDetails.viewCount` | `microformat.playerMicroformatRenderer.likeCount` when present | Video metadata |
| `/next` | Watch-page count fields | Like-button data, sometimes associated entities | Watch-page UI data |

`/updated_metadata` is the targeted endpoint for counter updates. `/player` also
provides useful counters. `/next` can expose display values, but a numeric field
can conflict with displayed text, including a zero placeholder on a Shorts watch
page. Do not silently replace a populated count with a contradictory placeholder.

### Request

**Endpoint:** `/updated_metadata`

```json
{
  "videoId": "VIDEO_ID"
}
```

For an update continuation, include both fields:

```json
{
  "videoId": "VIDEO_ID",
  "continuation": "METADATA_CONTINUATION_TOKEN"
}
```

The endpoint can also return changes to date, title, or description. Its response
is not a flat object containing only `viewCount` and `likeCount`.
[Metadata update requests](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/LiveChat.ts)

### Views

Inspect each `actions[]` entry for:

```text
updateViewershipAction.viewCount.videoViewCountRenderer
```

| Relative field | Meaning |
|---|---|
| `originalViewCount` | Numeric-string count when available |
| `unlabeledViewCountValue.simpleText` | Full displayed count, potentially with grouping separators |
| `viewCount.simpleText` or `viewCount.runs[]` | Count with its display label |
| `shortViewCount` / `extraShortViewCount` | Abbreviated count, such as `31K` |
| `isLive` | Live-viewership indicator |

Live viewership can describe concurrent viewers rather than cumulative video
views. Keep these meanings separate. An abbreviated display count is not an exact
integer. [View-count fields](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/VideoViewCount.ts)

### Likes

Inspect:

```text
frameworkUpdates.entityBatchUpdate.mutations[].payload.likeCountEntity
```

| Relative field | Meaning |
|---|---|
| `likeCountIfIndifferentNumber` | Numeric-string count for a viewer who has not liked the video |
| `likeCountIfLikedNumber` | Count for the liked state; can include a prospective extra like |
| `likeCountIfDislikedNumber` | Count for the disliked state |
| `expandedLikeCountIfIndifferent.content` | Expanded display count for the neutral state |
| `likeCountIfIndifferent.content` | Display count, possibly abbreviated |

**For anonymous requests, use `likeCountIfIndifferentNumber` when present.**
Do not use `likeCountIfLikedNumber` just because it is available; its state-dependent
value can be one higher. An authenticated viewer's current reaction must be considered
when interpreting state-dependent counts.

Some response versions instead use `updateToggleButtonTextAction`. Identify the
like control through `buttonId` and treat its text according to its precision.
[Button-text update fields](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/livechat/UpdateToggleButtonTextAction.ts)

### Update timing

```text
continuation.timedContinuationData.continuation
continuation.timedContinuationData.timeoutMs
```

`timeoutMs` is the server-provided wait before following that update continuation.
It is not a universal polling interval. A response may omit a counter that has no
available update; absence is not zero.

### Python: extract numeric counter updates

This function accepts a decoded `/updated_metadata` JSON object from an anonymous
request. It extracts numeric values only, preserving missing values as `None`.

```python
# example: metrics_parser

def numeric_count(value):
    if isinstance(value, bool):
        return None
    if isinstance(value, int):
        return value if value >= 0 else None
    if isinstance(value, str):
        value = value.strip()
        if value.isascii() and value.isdecimal():
            return int(value)
    return None


def extract_metrics(data):
    result = {
        "view_count": None,
        "live_viewer_count": None,
        "like_count": None,
    }

    for action in data.get("actions", []):
        renderer = (
            action.get("updateViewershipAction", {})
            .get("viewCount", {})
            .get("videoViewCountRenderer", {})
        )
        if renderer:
            field = "live_viewer_count" if renderer.get("isLive") else "view_count"
            result[field] = numeric_count(renderer.get("originalViewCount"))

    mutations = (
        data.get("frameworkUpdates", {})
        .get("entityBatchUpdate", {})
        .get("mutations", [])
    )
    for mutation in mutations:
        entity = mutation.get("payload", {}).get("likeCountEntity", {})
        if "likeCountIfIndifferentNumber" in entity:
            result["like_count"] = numeric_count(entity["likeCountIfIndifferentNumber"])

    timing = data.get("continuation", {}).get("timedContinuationData", {})
    result["continuation"] = timing.get("continuation")
    result["poll_after_ms"] = timing.get("timeoutMs")
    return result
```

The function targets this endpoint's numeric response shape. It does not turn
abbreviated labels into invented exact counts or parse unrelated watch-page data.

## Read top-level comments

**Endpoint:** `/next`

A video-ID request returns watch-page information and the comment entry point:

```json
{
  "videoId": "VIDEO_ID"
}
```

The same endpoint also returns recommendations and playlist context. Read the
requested video's comment section, not a continuation for the recommendation feed.
[Watch-page structure](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/VideoInfo.ts)

### 1. Obtain the comment-section token

Within:

```text
contents.twoColumnWatchNextResults.results.results.contents[]
```

Find the `itemSectionRenderer` whose `targetId` is `comments-section`. Its direct
items can contain:

```text
contents[].continuationItemRenderer.continuationEndpoint
  .continuationCommand.token
```

Send that token to `/next`:

```json
{
  "continuation": "COMMENT_SECTION_TOKEN"
}
```

A separate engagement-panel comments section may also exist. Use one consistent
comment surface rather than mixing tokens from different surfaces.

### 2. Select the ordering

The initial comment response can contain:

```text
commentsHeaderRenderer.sortMenu.sortFilterSubMenuRenderer.subMenuItems[]
```

Each item describes a sort order through fields such as `title` and `selected`.
For English WEB responses, the labels include `Top` and `Newest`. Select the desired
item and send its token:

```text
serviceEndpoint.continuationCommand.token
```

```json
{
  "continuation": "NEWEST_SORT_TOKEN"
}
```

Do not construct the token by changing characters. This API shape does not provide
a verified `since_comment_id` filter. Newest ordering and stable comment IDs support
identifying newly encountered comments, but ordering alone is not a completeness
guarantee. [Sorting and comment continuations](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/Comments.ts)

### 3. Read the top-level threads

Comment responses use commands such as:

```text
onResponseReceivedEndpoints[].reloadContinuationItemsCommand
onResponseReceivedEndpoints[].appendContinuationItemsAction
```

Select commands with `targetId == "comments-section"` and combine their
`continuationItems[]`. A response can separate the header and threads into different
commands with the same target. Do not assume the threads are always in array index 0.

A top-level item is a `commentThreadRenderer`. Its current comment view is:

```text
commentThreadRenderer.commentViewModel.commentViewModel
```

The view identifies the comment with `commentId` and references its data through
`commentKey`. The comment text and author can live elsewhere in the same response:

```text
frameworkUpdates.entityBatchUpdate.mutations[]
```

Build an index by `entityKey`, look up the view's `commentKey`, and read that
mutation's `payload.commentEntityPayload`.

| Field relative to `commentEntityPayload` | Information |
|---|---|
| `properties.commentId` | Stable comment ID |
| `properties.content.content` | Comment text |
| `properties.publishedTime` | Displayed publication time, often relative |
| `properties.replyLevel` | Reply depth; `0` identifies a top-level comment when supplied |
| `author.channelId` | Author's channel ID |
| `author.displayName` | Displayed author name or handle |
| `author.avatarThumbnailUrl` | Avatar URL |
| `author.isCreator` | Whether the author is the video's creator |
| `toolbar.likeCountNotliked` | Displayed comment likes in the neutral state; may be blank |
| `toolbar.likeCountA11y` | Accessible like-count label |
| `toolbar.replyCount` | Displayed reply count |
| `toolbar.replyCountA11y` | Accessible reply-count label |

A blank display count is not independently an exact numeric zero. Publication
labels such as `2 hours ago` do not provide an exact timestamp. Comment likes are
separate from video likes. [Comment entity structure](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/comments/CommentView.ts)

### 4. Fetch the next main page

Follow a `continuationItemRenderer` that is a **direct item of the main comment
list**. Its `continuationEndpoint.continuationCommand.token` retrieves another page.

Tokens nested under `commentRepliesRenderer` retrieve replies. Skip those for a
top-level-only result. Also avoid collecting nested reply previews as main comments.

Pinning, moderation, late visibility, and concurrent changes can affect ordering.
Encountering one stored comment ID is not a reliable universal stopping condition.
The displayed comment total can include content beyond the returned top-level
threads; it is not a pagination-completion test.

### Python: join comment views with their entities

The function below accepts a decoded WEB comment response. It handles the current
entity-based layout, returns top-level comments and the next main-page token, and
raises on an unresolved comment entity instead of silently dropping the record.

```python
# example: comments_parser

def extract_comments(data):
    mutations = (
        data.get("frameworkUpdates", {})
        .get("entityBatchUpdate", {})
        .get("mutations", [])
    )
    entities = {
        item["entityKey"]: item.get("payload", {})
        for item in mutations
        if "entityKey" in item
    }

    items = []
    for root in ("onResponseReceivedEndpoints", "onResponseReceivedActions", "onResponseReceivedCommands"):
        for entry in data.get(root, []):
            for kind in ("reloadContinuationItemsCommand", "appendContinuationItemsAction"):
                command = entry.get(kind, {})
                if command.get("targetId") == "comments-section":
                    items.extend(command.get("continuationItems", []))

    comments = []
    next_token = None
    for item in items:
        continuation = item.get("continuationItemRenderer", {})
        token = (
            continuation.get("continuationEndpoint", {})
            .get("continuationCommand", {})
            .get("token")
        )
        if token:
            next_token = token

        thread = item.get("commentThreadRenderer")
        if thread is None:
            continue
        view = thread.get("commentViewModel", {}).get("commentViewModel", {})
        entity = entities.get(view.get("commentKey"), {}).get("commentEntityPayload")
        if entity is None:
            raise ValueError("Unresolved comment entity or unsupported comment layout")

        properties = entity["properties"]
        if properties.get("replyLevel", 0) != 0:
            continue
        if properties["commentId"] != view.get("commentId"):
            raise ValueError("Comment ID and entity do not match")

        author = entity.get("author", {})
        comments.append({
            "comment_id": properties["commentId"],
            "text": properties["content"]["content"],
            "published_time_text": properties.get("publishedTime"),
            "author_channel_id": author.get("channelId"),
            "author_name": author.get("displayName"),
        })

    return comments, next_token
```

An empty result needs interpretation alongside the section/header and any errors.
It does not by itself distinguish disabled comments, a genuinely empty list, or
an unsupported response layout. Older clients can use `commentRenderer` objects
with inline fields instead of this entity-based layout.

## Search

**Endpoint:** `/search`

```json
{
  "query": "RizaNova"
}
```

Optional `params` can specify supported filters such as result type, upload date,
duration, and features. Use parameters provided by YouTube or a compatible encoder.
[Search filters](https://github.com/LuanRT/YouTube.js/blob/main/src/Innertube.ts)

Initial results commonly appear inside:

```text
contents.twoColumnSearchResultsRenderer.primaryContents.sectionListRenderer.contents[]
```

Results can include `videoRenderer`, `channelRenderer`, `lockupViewModel`, and
`shortsLockupViewModel`, including items inside shelves. Search continuations use
`{"continuation": "SEARCH_TOKEN"}` and can return items through
`onResponseReceivedCommands[].appendContinuationItemsAction`.

Search ranking, filters, and result limits mean this is not a complete channel
inventory or a guaranteed search through all spoken content. `estimatedResults`
is not the number of records guaranteed to be retrievable.

## Read a transcript

**Endpoint:** `/get_transcript`

Obtain the transcript request from the video's `/next` response. Find an
`engagementPanelSectionListRenderer` whose `panelIdentifier` is
`engagement-panel-searchable-transcript`, then inspect:

```text
content.continuationItemRenderer.continuationEndpoint.getTranscriptEndpoint.params
```

```json
{
  "params": "TRANSCRIPT_PARAMS_FROM_RESPONSE"
}
```

A successful response can contain `transcriptSegmentRenderer` entries with:

| Field | Information |
|---|---|
| `snippet` | Segment text, using YouTube's text representation |
| `startMs` | Segment start in milliseconds |
| `endMs` | Segment end in milliseconds |
| `startTimeText` | Displayed start time |

Follow language choices advertised by the transcript's language menu; each can
supply its own request parameters. [Transcript segments](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/TranscriptSegment.ts),
[language selection](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/TranscriptInfo.ts)

A transcript panel does not guarantee that its request will be accepted. A
`400 FAILED_PRECONDITION` response means request/session prerequisites were not
satisfied; it is not proof that the video has no transcript. Preserve this as an
access or request failure. This endpoint retrieves available text and does not
create a transcription for a video without one.

## Shorts viewer endpoints

### /reel/reel_item_watch

Obtain a `reelWatchEndpoint` from a Shorts card. Construct the read request using
its `videoId`, `playerParams`, and `params`:

```json
{
  "playerRequest": {
    "videoId": "SHORT_VIDEO_ID",
    "params": "PLAYER_PARAMS_FROM_CARD"
  },
  "params": "PARAMS_FROM_CARD",
  "disablePlayerResponse": false
}
```

The response can contain `status`, `overlay`, `engagementPanels`, and
`playerResponse`. When present, `playerResponse.videoDetails` and its microformat
can be interpreted like a `/player` response. Optional player parameters should
be omitted when not provided. [Request structure](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/endpoints/ReelWatchEndpoint.ts)

### /reel/reel_watch_sequence

```json
{
  "sequenceParams": "SEQUENCE_PARAMS_FROM_REEL_ENDPOINT"
}
```

Read `entries[]`, whose `command` describes the next viewing destination. A
`continuationEndpoint.continuationCommand.token` can be passed as the next
`sequenceParams`. This retrieves a Shorts viewing/recommendation sequence, not
all Shorts belonging to a specified channel. [Sequence handling](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/ytshorts/ShortFormVideoInfo.ts)

## Live chat

| Endpoint | Information |
|---|---|
| `/live_chat/get_live_chat` | Live chat events |
| `/live_chat/get_live_chat_replay` | Available archived chat events |

Start with the chat continuation advertised by the video's watch response:

```json
{
  "continuation": "CHAT_CONTINUATION_TOKEN"
}
```

The response can use `continuationContents.liveChatContinuation`, with `actions[]`
and another continuation. Follow the returned polling interval when supplied.
Events can include text, paid messages, membership events, replacements, and
removal notices. Replay events can carry timing relative to video playback.

Live chat is separate from comments below a video. Replay availability depends
on the video and the chat settings. [Live-chat handling](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/LiveChat.ts)

## Other reading endpoints

These expose other YouTube surfaces. They do not replace uploads, metadata,
counter, or comment requests.

| Endpoint | Information | Request/context |
|---|---|---|
| `/guide` | Sidebar navigation destinations | Client/account context. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/Innertube.ts) |
| `/get_panel` | A panel advertised by the interface | Returned `panelId` and `params`; follow the relevant panel command. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/endpoints/ShowEngagementPanelEndpoint.ts) |
| `/share/get_share_panel` | Sharing-panel information | `serializedSharedEntity` and optional `clientParams` from the returned share command. Fetching the panel does not send a message. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/endpoints/ShareEntityServiceEndpoint.ts) |
| `/music/get_search_suggestions` | YouTube Music autocomplete suggestions | `input` and a compatible Music client configuration. [Source](https://github.com/tombulled/innertube/blob/main/innertube/clients.py) |
| `/music/get_queue` | YouTube Music queue information | `videoIds` or `playlistId` and compatible Music context. [Source](https://github.com/tombulled/innertube/blob/main/innertube/clients.py) |
| `/account/accounts_list` | Identities available to the signed-in user | Authenticated account session and account-list request configuration. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/core/managers/AccountManager.ts) |
| `/notification/get_notification_menu` | The signed-in account's notifications | Authenticated session; inbox request type and returned continuations. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/youtube/NotificationsMenu.ts) |
| `/notification/get_unseen_count` | Unseen notification count | Authenticated session. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/Innertube.ts) |
| `/playlist/get_add_to_playlist` | Playlist choices for the Save dialog | Authenticated session, video/playlist identifiers, and relevant returned parameters; adding is a separate action. [Source](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/endpoints/AddToPlaylistServiceEndpoint.ts) |

Music requests need a compatible Music client and host. Account information is
limited to the authenticated account. Notifications are not a complete upload
feed for arbitrary channels.

## Pagination and response parsing

### Collection-specific tokens

A response can advertise several continuations at once: playlist pages, comments,
replies, recommendations, or metadata updates. Follow the token belonging to the
intended collection and the endpoint named by its command metadata.

Common response envelopes include:

```text
onResponseReceivedActions[]
onResponseReceivedEndpoints[]
onResponseReceivedCommands[]
```

Common commands inside them include `appendContinuationItemsAction` and
`reloadContinuationItemsCommand`, with items in `continuationItems[]`. Preserve
`targetId` so that unrelated collections are not combined.
[Continuation routing](https://github.com/LuanRT/YouTube.js/blob/main/src/parser/classes/commands/ContinuationCommand.ts)

### Text formats

YouTube text is commonly represented as:

```json
{
  "simpleText": "Example text"
}
```

```json
{
  "runs": [
    {"text": "First part "},
    {"text": "second part"}
  ]
}
```

```json
{
  "content": "Example attributed text"
}
```

For `runs`, concatenate each run's `text` in order. Preserve Unicode, whitespace,
and line breaks. Attributed text can also carry formatting or attachment metadata.

### Parsing rules

- Scope extraction to the requested content area before walking nested objects.
- Join entity references using their keys; do not interpret entity keys as video IDs.
- Deduplicate records by stable IDs, including across pages.
- Detect repeated tokens or repeated pages to prevent infinite pagination.
- Do not interpret a transport, access, or parsing failure as an empty collection.
- Treat an unknown renderer layout as a parser compatibility issue.

## Errors and data semantics

| Response or condition | Meaning and handling |
|---|---|
| HTTP 200 | Transport success only; inspect the requested data and application status |
| JSON `error` object | API-level error; preserve its code, status, and message |
| `playabilityStatus.status == "ERROR"` with no video details | The request did not provide usable video metadata |
| `LOGIN_REQUIRED` with a sign-in or automated-access message | Access/session restriction; do not infer that the video was deleted |
| `UNPLAYABLE` with populated video details | Playback is restricted, but metadata may remain usable |
| HTTP 400 / `FAILED_PRECONDITION` | Invalid request or unmet prerequisites; do not treat as a successful empty result |
| HTTP 401 or 403 | Authentication or access failure |
| HTTP 429 | Rate limiting; honor `Retry-After` when supplied and back off |
| Timeout or HTTP 5xx | Potential transient failure; use bounded retries |
| HTML or invalid JSON | Unexpected response, redirect, or access page; not endpoint data |
| Missing views, likes, or dates | Unknown/unavailable, not automatically zero or false |
| A count decreases | Possible cache differences, reconciliation, moderation, or removed reactions; do not assume counts are monotonic |
| Date-only or relative publication text | Lower precision than a timestamp; preserve that distinction |
| No comment continuation | End of the returned collection, not proof that all historical comments are visible |

Counter responses are observations, not an exact event history. Fields from
separate endpoints can have different cache ages. A fresh request does not guarantee
that every count was calculated at that exact instant.

Unavailable fields, hidden likes, disabled comments, deleted content, access
restrictions, and changed response layouts are distinct states. Requests can expose
only the content made available to the chosen client/session. No fixed public
InnerTube requests-per-second allowance or universal page size is specified here.
