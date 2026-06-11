package response

import (
	"encoding/json"
	"testing"
)

func TestBranchUnmarshalMatchesLiveListResponseShape(t *testing.T) {
	payload := []byte(`{
		"name": "master",
		"commit": {
			"sha": "e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a",
			"url": "https://gitee.com/api/v5/repos/yimeng555/example/commits/e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a",
			"commit": {
				"author": {
					"name": "yimeng555",
					"date": "2026-06-11T20:29:54+08:00",
					"email": "17183222+yimeng555@user.noreply.gitee.com"
				},
				"committer": {
					"name": "Gitee",
					"date": "2026-06-11T20:29:54+08:00",
					"email": "noreply@gitee.com"
				},
				"message": "Initial commit"
			}
		},
		"protected": false,
		"protection_url": "https://gitee.com/api/v5/repos/yimeng555/example/branches/master/protection"
	}`)

	var branch Branch
	if err := json.Unmarshal(payload, &branch); err != nil {
		t.Fatalf("json.Unmarshal(branch) returned error: %v", err)
	}

	if branch.Commit.SHA != "e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a" {
		t.Fatalf("branch.Commit.SHA = %q, want %q", branch.Commit.SHA, "e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a")
	}
}

func TestCompleteBranchUnmarshalMatchesLiveDetailResponseShape(t *testing.T) {
	payload := []byte(`{
		"name": "master",
		"commit": {
			"sha": "e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a",
			"url": "https://gitee.com/api/v5/repos/yimeng555/example/commits/e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a",
			"commit": {
				"author": {
					"name": "yimeng555",
					"date": "2026-06-11T12:29:54+00:00",
					"email": "17183222+yimeng555@user.noreply.gitee.com"
				},
				"url": "https://gitee.com/api/v5/repos/yimeng555/example/git/commits/e80c3720a55b73df7a3f9cd3de35a8ec6fe9976a",
				"message": "Initial commit",
				"tree": {
					"sha": "89c99eab704e112e2ef512a05f1408b72990c261",
					"url": "https://gitee.com/api/v5/repos/yimeng555/example/git/trees/89c99eab704e112e2ef512a05f1408b72990c261"
				},
				"committer": {
					"name": "Gitee",
					"date": "2026-06-11T12:29:54+00:00",
					"email": "noreply@gitee.com"
				}
			},
			"author": {
				"avatar_url": "https://foruda.gitee.com/avatar.png",
				"url": "https://gitee.com/api/v5/users/yimeng555",
				"id": 17183222,
				"login": "yimeng555"
			},
			"parents": [],
			"committer": {
				"avatar_url": "https://foruda.gitee.com/avatar2.png",
				"url": "https://gitee.com/api/v5/users/gitee-bot",
				"id": 10186697,
				"login": "gitee-bot"
			}
		},
		"_links": {
			"html": "https://gitee.com/yimeng555/example/tree/master",
			"self": "https://gitee.com/api/v5/repos/yimeng555/example/branches/master"
		},
		"protected": false,
		"protection_url": "https://gitee.com/api/v5/repos/yimeng555/example/branches/master/protection"
	}`)

	var branch CompleteBranch
	if err := json.Unmarshal(payload, &branch); err != nil {
		t.Fatalf("json.Unmarshal(complete branch) returned error: %v", err)
	}

	if branch.Links.HTML == "" {
		t.Fatal("branch.Links.HTML is empty")
	}
}
