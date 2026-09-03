package tunnetcontrol

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func mapperFixtures(id string) (json.RawMessage, json.RawMessage) {
	bootstrap := json.RawMessage(`{"schema_version":2,"access":{"state":"ready"},"runtime":{"client_id":"` + id + `","entry_nodes":[{"name":"entry-a","ipv4":["203.0.113.10"],"front_proxy":{"endpoint":"http://127.0.0.1:8080","headers":{"Host":"front.example","X-T5-Auth":"auth"}}}],"hosts":[{"slug":"sin-03","online":true,"vless_encryption_key":"key"}]}}`)
	sync := json.RawMessage(`{"schema_version":2,"access":{"state":"ready"},"runtime":{"client_id":"` + id + `","entry_nodes":[{"name":"entry-a","ipv4":["203.0.113.10"],"front_proxy":{"endpoint":"http://127.0.0.1:8080","headers":{"Host":"front.example","X-T5-Auth":"auth"}}}],"hosts":[{"slug":"sin-03","online":true,"vless_encryption_key":"key"}]}}`)
	return bootstrap, sync
}

func TestMapSnapshotBuildsValidatedPublication(t *testing.T) {
	identity := testIdentity(t)
	bootstrap, synced := mapperFixtures(identity.ClientID)
	now := time.Now()
	snapshot, err := MapSnapshotWithOptions(context.Background(), bootstrap, synced, ECHConfig{QueryName: echQueryName, ConfigList: "AQID", ExpiresAt: now.Add(time.Hour).UnixMilli()}, identity, MapOptions{
		DataAuthority:     "sin-03.data.example",
		PreviousEntryNode: "entry-a",
		RouteProber:       fixtureProber{"203.0.113.10": {latency: time.Millisecond}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.ActiveEntryNode != "entry-a" || snapshot.SelectedHost != "sin-03" || snapshot.XHTTPAuthority != "sin-03.data.example" || snapshot.ValidatedRoute != "203.0.113.10" || !strings.HasPrefix(snapshot.XHTTPPath, "/api/v1/sync/") {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	if err := validateSnapshot(snapshot, now); err != nil {
		t.Fatal(err)
	}
}

func TestMapSnapshotFailsClosedOnUnknownSchema(t *testing.T) {
	identity := testIdentity(t)
	_, err := MapSnapshotWithOptions(context.Background(), json.RawMessage(`{}`), json.RawMessage(`{"schema_version":3,"access":{"state":"ready"},"runtime":{}}`), ECHConfig{}, identity, MapOptions{})
	if err == nil || !strings.Contains(err.Error(), "schema/runtime") {
		t.Fatalf("expected schema rejection, got %v", err)
	}
}

func TestMapSnapshotFailsClosedOnInvalidSelectionOrRoute(t *testing.T) {
	identity := testIdentity(t)
	bootstrap, synced := mapperFixtures(identity.ClientID)
	ech := ECHConfig{QueryName: echQueryName, ConfigList: "AQID", ExpiresAt: time.Now().Add(time.Hour).UnixMilli()}
	for _, options := range []MapOptions{
		{DataAuthority: "invalid", PreviousEntryNode: "entry-a", RouteProber: fixtureProber{"203.0.113.10": {latency: time.Millisecond}}},
		{DataAuthority: "sin-03.data.example", PreviousEntryNode: "missing", RouteProber: fixtureProber{"203.0.113.10": {latency: time.Millisecond}}},
		{DataAuthority: "sin-03.data.example", PreviousEntryNode: "entry-a", RouteProber: fixtureProber{}},
	} {
		if _, err := MapSnapshotWithOptions(context.Background(), bootstrap, synced, ech, identity, options); err == nil {
			t.Fatalf("accepted invalid options: %#v", options)
		}
	}
}
