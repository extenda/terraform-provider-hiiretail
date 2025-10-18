package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExerciseTestHelpers(t *testing.T) {
	// makeResp
	resp := makeResp(200, `{"ok":true}`)
	require.NotNil(t, resp)

	// seqRoundTripper: supply responses and call RoundTrip directly
	seq := &seqRoundTripper{responses: []*http.Response{makeResp(200, `ok`)}}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.test", nil)
	require.NoError(t, err)
	r, err := seq.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, r)

	// captureRT: ensure RoundTrip stores the request
	capt := &captureRT{resp: makeResp(201, `{"created":true}`)}
	r2, err := capt.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, r2)
	require.NotNil(t, capt.req)
}
