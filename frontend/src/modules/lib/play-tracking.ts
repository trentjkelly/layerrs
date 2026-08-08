import { get } from 'svelte/store';
import { audio, playedTracksThisSession } from '../../stores/player';
import { recordTrackPlay } from '../requests/track-requests';
import { logger } from './logger';

const PLAY_THRESHOLD_SECONDS = 10;
const SEEK_DETECTION_THRESHOLD = 1.5;

export function trackPlayProgress(trackId: number, urlBase: string): () => void {
    const $audio = get(audio);
    if (!$audio) {
        return () => {};
    }

    let accumulatedTime = 0;
    let lastCurrentTime = 0;
    let hasRecorded = false;
    let lastTimeupdate = 0;

    const session = get(playedTracksThisSession);
    if (session.has(trackId)) {
        hasRecorded = true;
    }

    const handleTimeUpdate = () => {
        const now = performance.now();
        const currentTime = $audio.currentTime;

        if (lastTimeupdate !== 0 && !$audio.paused && !$audio.seeking) {
            const wallDelta = (now - lastTimeupdate) / 1000;
            const timeDelta = currentTime - lastCurrentTime;

            // Only count time if the playhead moved forward roughly in real-time
            // (ignore seeks, buffering, and backward jumps)
            if (timeDelta > 0 && timeDelta < SEEK_DETECTION_THRESHOLD && wallDelta < SEEK_DETECTION_THRESHOLD * 2) {
                accumulatedTime += Math.min(timeDelta, wallDelta);
            }
        }

        lastCurrentTime = currentTime;
        lastTimeupdate = now;

        if (!hasRecorded && accumulatedTime >= PLAY_THRESHOLD_SECONDS) {
            recordPlay();
        }
    };

    const handleEnded = () => {
        if (!hasRecorded) {
            recordPlay();
        }
    };

    const recordPlay = () => {
        const session = get(playedTracksThisSession);
        if (session.has(trackId)) {
            hasRecorded = true;
            return;
        }

        hasRecorded = true;
        playedTracksThisSession.update((set) => {
            set.add(trackId);
            return set;
        });

        recordTrackPlay(urlBase, trackId).catch((error) => {
            logger.error(`Failed to record track play: ${error}`);
            // Allow a retry on next playback by resetting the flag
            playedTracksThisSession.update((set) => {
                set.delete(trackId);
                return set;
            });
            hasRecorded = false;
        });
    };

    $audio.addEventListener('timeupdate', handleTimeUpdate);
    $audio.addEventListener('ended', handleEnded);

    return () => {
        $audio.removeEventListener('timeupdate', handleTimeUpdate);
        $audio.removeEventListener('ended', handleEnded);
    };
}
