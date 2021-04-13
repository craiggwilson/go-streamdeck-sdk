package streamdeck

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"
)

var (
	_ ActionPublisher         = &coreActionPublisher{}
	_ ActionInstancePublisher = &coreActionInstancePublisher{}
)

// Publisher publishes events for a plugin. It is an alias for a Publisher.
type Publisher = streamdeckcore.Publisher

// ActionPublisher publishes events for an Action, filling in details specific to the Action.
type ActionPublisher interface {
	Publisher

	GetGlobalSettings() error
	GetSettings(eventContext EventContext) error
	LogMessage(payload streamdeckevent.LogMessagePayload) error
	OpenURL(payload streamdeckevent.OpenURLPayload) error
	SendToPropertyInspector(eventContext EventContext, payload json.RawMessage) error
	SetGlobalSettings(settings json.RawMessage) error
	SetImage(eventContext EventContext, payload streamdeckevent.SetImagePayload) error
	SetSettings(eventContext EventContext, settings json.RawMessage) error
	SetState(eventContext EventContext, payload streamdeckevent.SetStatePayload) error
	SetTitle(eventContext EventContext, payload streamdeckevent.SetTitlePayload) error
	ShowAlert(eventContext EventContext) error
	ShowOK(eventContext EventContext) error
	SwitchToProfile(eventContext EventContext, payload streamdeckevent.SwitchToProfilePayload) error
}

func newCoreActionPublisher(
	pluginUUID PluginUUID,
	actionUUID ActionUUID,
	corePublisher Publisher) *coreActionPublisher {

	return &coreActionPublisher{
		pluginUUID:    pluginUUID,
		actionUUID:    actionUUID,
		corePublisher: corePublisher,
	}
}

type coreActionPublisher struct {
	deviceUUID    DeviceUUID
	pluginUUID    PluginUUID
	actionUUID    ActionUUID
	corePublisher Publisher
}

func (p *coreActionPublisher) GetGlobalSettings() error {
	event := streamdeckevent.GetGlobalSettings{
		Event:   streamdeckevent.GetGlobalSettingsName,
		Context: p.pluginUUID,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) GetSettings(eventContext EventContext) error {
	event := streamdeckevent.GetSettings{
		Event:   streamdeckevent.GetSettingsName,
		Context: eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) LogMessage(payload streamdeckevent.LogMessagePayload) error {
	event := streamdeckevent.LogMessage{
		Event:   streamdeckevent.LogMessageName,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) OpenURL(payload streamdeckevent.OpenURLPayload) error {
	event := streamdeckevent.OpenURL{
		Event:   streamdeckevent.OpenURLName,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SendToPropertyInspector(eventContext EventContext, payload json.RawMessage) error {
	event := streamdeckevent.SendToPropertyInspector{
		Action:  p.actionUUID,
		Event:   streamdeckevent.SendToPropertyInspectorName,
		Context: eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SetGlobalSettings(settings json.RawMessage) error {
	event := streamdeckevent.SetGlobalSettings{
		Event:   streamdeckevent.SetGlobalSettingsName,
		Context: p.pluginUUID,
		Payload: settings,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SetImage(eventContext EventContext, payload streamdeckevent.SetImagePayload) error {
	event := streamdeckevent.SetImage{
		Event:   streamdeckevent.SetImageName,
		Context: eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SetSettings(eventContext EventContext, settings json.RawMessage) error {
	event := streamdeckevent.SetSettings{
		Event:   streamdeckevent.SetSettingsName,
		Context: eventContext,
		Payload: settings,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SetState(eventContext EventContext, payload streamdeckevent.SetStatePayload) error {
	event := streamdeckevent.SetState{
		Event:   streamdeckevent.SetStateName,
		Context: eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SetTitle(eventContext EventContext, payload streamdeckevent.SetTitlePayload) error {
	event := streamdeckevent.SetTitle{
		Event:   streamdeckevent.SetTitleName,
		Context: eventContext,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) ShowAlert(eventContext EventContext) error {
	event := streamdeckevent.ShowAlert{
		Event:   streamdeckevent.ShowAlertName,
		Context: eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) ShowOK(eventContext EventContext) error {
	event := streamdeckevent.ShowOK{
		Event:   streamdeckevent.ShowOKName,
		Context: eventContext,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) SwitchToProfile(eventContext EventContext, payload streamdeckevent.SwitchToProfilePayload) error {
	event := streamdeckevent.SwitchToProfile{
		Event:   streamdeckevent.SwitchToProfileName,
		Context: eventContext,
		Device:  p.deviceUUID,
		Payload: payload,
	}

	return p.publish(event.Event, event)
}

func (p *coreActionPublisher) PublishEvent(raw json.RawMessage) error {
	return p.corePublisher.PublishEvent(raw)
}

func (p *coreActionPublisher) publish(eventName EventName, event interface{}) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event %q: %w", eventName, err)
	}

	if err = p.PublishEvent(raw); err != nil {
		return fmt.Errorf("publishing event %q: %w", eventName, err)
	}

	return nil
}

// ActionInstancePublisher publishes events for an ActionInstance, filling in details specific to the ActionInstance.
type ActionInstancePublisher interface {
	Publisher

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
	eventContext EventContext,
	corePublisher ActionPublisher) *coreActionInstancePublisher {

	return &coreActionInstancePublisher{
		eventContext:  eventContext,
		actionPublisher: corePublisher,
	}
}

type coreActionInstancePublisher struct {
	eventContext  EventContext
	actionPublisher ActionPublisher
}

func (p *coreActionInstancePublisher) GetGlobalSettings() error {
	return p.actionPublisher.GetGlobalSettings()
}

func (p *coreActionInstancePublisher) GetSettings() error {
	return p.actionPublisher.GetSettings(p.eventContext)
}

func (p *coreActionInstancePublisher) LogMessage(payload streamdeckevent.LogMessagePayload) error {
	return p.actionPublisher.LogMessage(payload)
}

func (p *coreActionInstancePublisher) OpenURL(payload streamdeckevent.OpenURLPayload) error {
	return p.actionPublisher.OpenURL(payload)
}

func (p *coreActionInstancePublisher) SendToPropertyInspector(payload json.RawMessage) error {
	return p.actionPublisher.SendToPropertyInspector(p.eventContext, payload)
}

func (p *coreActionInstancePublisher) SetGlobalSettings(settings json.RawMessage) error {
	return p.actionPublisher.SetGlobalSettings(settings)
}

func (p *coreActionInstancePublisher) SetImage(payload streamdeckevent.SetImagePayload) error {
	return p.actionPublisher.SetImage(p.eventContext, payload)
}

func (p *coreActionInstancePublisher) SetSettings(settings json.RawMessage) error {
	return p.actionPublisher.SetSettings(p.eventContext, settings)
}

func (p *coreActionInstancePublisher) SetState(payload streamdeckevent.SetStatePayload) error {
	return p.actionPublisher.SetState(p.eventContext, payload)
}

func (p *coreActionInstancePublisher) SetTitle(payload streamdeckevent.SetTitlePayload) error {
	return p.actionPublisher.SetTitle(p.eventContext, payload)
}

func (p *coreActionInstancePublisher) ShowAlert() error {
	return p.actionPublisher.ShowAlert(p.eventContext)
}

func (p *coreActionInstancePublisher) ShowOK() error {
	return p.actionPublisher.ShowOK(p.eventContext)
}

func (p *coreActionInstancePublisher) SwitchToProfile(payload streamdeckevent.SwitchToProfilePayload) error {
	return p.actionPublisher.SwitchToProfile(p.eventContext, payload)
}

func (p *coreActionInstancePublisher) PublishEvent(raw json.RawMessage) error {
	return p.actionPublisher.PublishEvent(raw)
}
