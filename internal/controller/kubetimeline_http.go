package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	timelinev1alpha1 "github.com/Thinnesen/kubetimeline/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func StartHTTPServer(k8sClient client.Client) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
  <title>KubeTimeline Dashboard</title>
  <style>
    body { font-family: sans-serif; background: #f4f4f4; margin: 0; padding: 2em; }
    h1 { color: #2c3e50; }
    .event { background: #fff; margin: 1em 0; padding: 1em; border-radius: 8px; box-shadow: 0 2px 4px #0001; }
    .event time { color: #888; font-size: 0.9em; }
  </style>
</head>
<body>
  <h1>KubeTimeline Events</h1>
  <div id="timeline"></div>
  <script>
    fetch('/timeline')
      .then(res => res.json())
    																																																		  .then(data => {
        let html = '';
        (data.items || []).forEach(tl => {
          html += '<h2>' + tl.metadata.name + '</h2>';
          (tl.status && tl.status.events || []).forEach(ev => {
            html += '<div class="event">' + ev + '</div>';
          });
        });
        document.getElementById('timeline').innerHTML = html;
      })
      .catch(err => {
        document.getElementById('timeline').innerHTML = '<b>Error loading timeline:</b> ' + err;
      });
  </script>
</body>
</html>	
		`))
	})

	http.HandleFunc("/timeline", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		var timelines timelinev1alpha1.KubeTimelineList
		if err := k8sClient.List(ctx, &timelines); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(timelines)
	})
	port := "8080"
	if p := os.Getenv("HTTP_PORT"); p != "" {
		port = p
	}
	go http.ListenAndServe(":"+port, nil)
}
