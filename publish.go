package streamdeck

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"
)

// ActionInstancePublisher publishes events for an ActionInstance, filling in details specific to the ActionInstance.
type ActionInstancePublisher interface {
	streamdeckcore.Publisher

	GetGlobalSettings() error
	GetSettings() error
	LogMessage(payload streamdeckevent.LogMessagePayload) error
	OpenURL(payload streamdeckevent.OpenURLPayload) error
	SendToPropertyInspector(payload json.RawMessage) error
	SetGlobalSettings(settings json.RawMessage) error
	SetImage(payload streamdeckevent.SetImagePayload) error
	SetSettings(settings json.RawMessage) error
	SetState(payload streamdeckevent.SetStatePayload) error
	SetTitle(payload streamdeckevent.SetTitlePayload) error
	ShowAlert() error
	ShowOK() error
	SwitchToProfile(payload streamdeckevent.SwitchToProfilePayload) error
}

func newCoreActionInstancePublisher(
	pluginUUID streamdeckcore.PluginUUID,
	actionUUID    streamdeckcore.ActionUUID,
	eventContext  streamdeckcore.EventContext,
	corePublisher streamdeckcore.Publisher) *coreActionInstancePublisher {

	return &coreActionInstancePublisher{
		pluginUUID: pluginUUID,
		actionUUID: actionUUID,
		eventContext: eventContext,
		corePublisher: corePublisher,
	}
}

type coreActionInstancePublisher struct {
	deviceUUID streamdeckcore.DeviceUUID
	pluginUUID streamdeckcore.PluginUUID
	actionUUID    streamdeckcore.ActionUUID
	eventContext  streamdeckcore.EventContext
	corePublisher streamdeckcore.Publisher
}

func (p *coreActionInstancePublisher) GetGlobalSettings() error {
	event := streamdeckevent.GetGlobalSettings{
		Event: streamdeckevent.GetGlobalSettingsName,
		Context: p.pluginUUID,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) GetSettings() error {
	event := streamdeckevent.GetSettings{
		Event: streamdeckevent.GetSettingsName,
		Context: p.eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) LogMessage(payload streamdeckevent.LogMessagePayload) error {
	event := streamdeckevent.LogMessage{
		Event: streamdeckevent.LogMessageName,
		Payload: payload,
	}
	
	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) OpenURL(payload streamdeckevent.OpenURLPayload) error {
	event := streamdeckevent.OpenURL{
		Event: streamdeckevent.OpenURLName,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SendToPropertyInspector(payload json.RawMessage) error {
	event := streamdeckevent.SendToPropertyInspector{
		Action: p.actionUUID,
		Event: streamdeckevent.SendToPropertyInspectorName,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SetGlobalSettings(settings json.RawMessage) error {
	event := streamdeckevent.SetGlobalSettings{
		Event: streamdeckevent.SetGlobalSettingsName,
		Context: p.pluginUUID,
		Payload: settings,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SetImage(payload streamdeckevent.SetImagePayload) error {
	event := streamdeckevent.SetImage{
		Event: streamdeckevent.SetImageName,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SetSettings(settings json.RawMessage) error {
	event := streamdeckevent.SetSettings{
		Event: streamdeckevent.SetSettingsName,
		Context: p.eventContext,
		Payload: settings,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SetState(payload streamdeckevent.SetStatePayload) error {
	event := streamdeckevent.SetState{
		Event: streamdeckevent.SetStateName,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SetTitle(payload streamdeckevent.SetTitlePayload) error {
	event := streamdeckevent.SetTitle{
		Event: streamdeckevent.SetTitleName,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) ShowAlert() error {
	event := streamdeckevent.ShowAlert{
		Event: streamdeckevent.ShowAlertName,
		Context: p.eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) ShowOK() error {
	event := streamdeckevent.ShowOK{
		Event: streamdeckevent.ShowOKName,
		Context: p.eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) SwitchToProfile(payload streamdeckevent.SwitchToProfilePayload) error {
	event := streamdeckevent.SwitchToProfile{
		Event: streamdeckevent.SwitchToProfileName,
		Context: p.eventContext,
		Device: p.deviceUUID,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionInstancePublisher) PublishEvent(raw json.RawMessage) error {
	return p.corePublisher.PublishEvent(raw)
}

func (p *coreActionInstancePublisher) publish(eventName streamdeckcore.EventName, event interface{}) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event %q: %w", eventName, err)
	}

	if err = p.PublishEvent(raw); err != nil {
		return fmt.Errorf("publishing event %q: %w", eventName, err)
	}

	return nil
}
