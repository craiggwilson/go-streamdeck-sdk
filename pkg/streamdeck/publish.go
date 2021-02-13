package streamdeck

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

type Publisher interface {
	GetGlobalSettings() error
	GetSettings() error
	SetGlobalSettings(settings json.RawMessage) error
	SetImage(event streamdeckcore.SetImagePayload) error
	SetSettings(event json.RawMessage) error
	SetState(event streamdeckcore.SetStatePayload) error
	SetTitle(event streamdeckcore.SetTitlePayload) error
	ShowAlert() error
	ShowOK() error

	PublishRaw(raw json.RawMessage) error
}

func newCorePublisherAdapter(
	pluginUUID streamdeckcore.PluginUUID,
	actionUUID    streamdeckcore.ActionUUID,
	eventContext  streamdeckcore.EventContext,
	corePublisher streamdeckcore.Publisher) *corePublisherAdapter {

	return &corePublisherAdapter{
		pluginUUID: pluginUUID,
		actionUUID: actionUUID,
		eventContext: eventContext,
		corePublisher: corePublisher,
	}
}

type corePublisherAdapter struct {
	pluginUUID streamdeckcore.PluginUUID
	actionUUID    streamdeckcore.ActionUUID
	eventContext  streamdeckcore.EventContext
	corePublisher streamdeckcore.Publisher
}

func (p *corePublisherAdapter) GetGlobalSettings() error {
	event := streamdeckcore.GetGlobalSettingsEvent{
		Event: streamdeckcore.GetGlobalSettings,
		Context: p.pluginUUID,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) GetSettings() error {
	event := streamdeckcore.GetSettingsEvent{
		Event: streamdeckcore.GetSettings,
		Context: p.eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) SetGlobalSettings(settings json.RawMessage) error {
	event := streamdeckcore.SetGlobalSettingsEvent{
		Event: streamdeckcore.SetGlobalSettings,
		Context: p.pluginUUID,
		Payload: settings,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) SetImage(payload streamdeckcore.SetImagePayload) error {
	event := streamdeckcore.SetImageEvent{
		Event: streamdeckcore.SetImage,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) SetSettings(settings json.RawMessage) error {
	event := streamdeckcore.SetSettingsEvent{
		Event: streamdeckcore.SetSettings,
		Context: p.eventContext,
		Payload: settings,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) SetState(payload streamdeckcore.SetStatePayload) error {
	event := streamdeckcore.SetStateEvent{
		Event: streamdeckcore.SetState,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) SetTitle(payload streamdeckcore.SetTitlePayload) error {
	event := streamdeckcore.SetTitleEvent{
		Event: streamdeckcore.SetTitle,
		Context: p.eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) ShowAlert() error {
	event := streamdeckcore.ShowAlertEvent{
		Event: streamdeckcore.ShowAlert,
		Context: p.eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) ShowOK() error {
	event := streamdeckcore.ShowOKEvent{
		Event: streamdeckcore.ShowOK,
		Context: p.eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *corePublisherAdapter) PublishRaw(raw json.RawMessage) error {
	return p.corePublisher.PublishEvent(raw)
}

func (p *corePublisherAdapter) publish(eventName streamdeckcore.EventName, event interface{}) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event %q: %w", eventName, err)
	}

	if err = p.PublishRaw(raw); err != nil {
		return fmt.Errorf("publishing event %q: %w", eventName, err)
	}

	return nil
}
