
// maybe metadata has enough timestamp data to stitch the data back
// maybe use channel and use sliding window to stitch the data
// use heap to keep the chunks sorted within sliding window
// keep track of the first item (earliest timestamp known so far)
// whenever new chunk comes in, determine whether to flush the chunks
// give reasonable flush timeout (5 min?) before all chunks that timed out are flushed
// this might introduce slight delay as we track the timed out chunks flush.
// whenever new chunk comes in, compare with the earliest timestamp and stitch
// when cleaning up after timer expired, flush rest of the chunks
// at this point, all chunks should be downloaded?