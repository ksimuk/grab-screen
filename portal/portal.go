package portal

import (
	"fmt"
	"net/url"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	portalDest          = "org.freedesktop.portal.Desktop"
	portalPath          = "/org/freedesktop/portal/desktop"
	screenshotInterface = "org.freedesktop.portal.Screenshot"
	requestInterface    = "org.freedesktop.portal.Request"
)

// TakeScreenshot requests a screenshot from the XDG Desktop Portal.
// It returns the URI of the saved screenshot.
func TakeScreenshot() (string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return "", fmt.Errorf("failed to connect to session bus: %w", err)
	}
	defer conn.Close()

	obj := conn.Object(portalDest, portalPath)

	// Options for the screenshot request
	// handle_token is useful for tracking the request, but we can let the portal generate one if we don't supply it.
	// interactive: true allows the user to choose options (like area selection).
	options := map[string]dbus.Variant{
		"interactive": dbus.MakeVariant(true),
	}

	// Call the Screenshot method
	// The signature is Screenshot(parent_window: s, options: a{sv}) -> (response_path: o)
	// parent_window can be empty string if we don't have a window.
	var responsePath dbus.ObjectPath
	err = obj.Call(screenshotInterface+".Screenshot", 0, "", options).Store(&responsePath)
	if err != nil {
		return "", fmt.Errorf("failed to call Screenshot method: %w", err)
	}

	// We need to listen for the Response signal on the returned request object path
	// The signal is org.freedesktop.portal.Request.Response(response: u, results: a{sv})
	// response: 0 = Success, 1 = User Cancelled, 2 = Failed

	// Add a match rule for the signal
	if err := conn.AddMatchSignal(
		dbus.WithMatchInterface(requestInterface),
		dbus.WithMatchMember("Response"),
		dbus.WithMatchObjectPath(responsePath),
	); err != nil {
		return "", fmt.Errorf("failed to add match signal: %w", err)
	}
	defer conn.RemoveMatchSignal(
		dbus.WithMatchInterface(requestInterface),
		dbus.WithMatchMember("Response"),
		dbus.WithMatchObjectPath(responsePath),
	)

	c := make(chan *dbus.Signal, 1)
	conn.Signal(c)

	// Wait for the signal
	select {
	case sig := <-c:
		if len(sig.Body) < 2 {
			return "", fmt.Errorf("unexpected signal body length")
		}

		responseCode, ok := sig.Body[0].(uint32)
		if !ok {
			return "", fmt.Errorf("unexpected type for response code")
		}

		if responseCode != 0 {
			return "", fmt.Errorf("screenshot request failed or cancelled (code: %d)", responseCode)
		}

		results, ok := sig.Body[1].(map[string]dbus.Variant)
		if !ok {
			return "", fmt.Errorf("unexpected type for results")
		}

		uriVariant, ok := results["uri"]
		if !ok {
			return "", fmt.Errorf("no uri in results")
		}

		uri, ok := uriVariant.Value().(string)
		if !ok {
			return "", fmt.Errorf("uri is not a string")
		}

		// The URI is usually in the format file:///path/to/file.png
		u, err := url.Parse(uri)
		if err != nil {
			return "", fmt.Errorf("failed to parse uri: %w", err)
		}

		if u.Scheme != "file" {
			return "", fmt.Errorf("unexpected scheme in uri: %s", u.Scheme)
		}

		return u.Path, nil
	case <-time.After(5 * time.Minute):
		return "", fmt.Errorf("timeout waiting for screenshot portal response")
	}
}
