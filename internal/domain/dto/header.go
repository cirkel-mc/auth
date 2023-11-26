package dto

type RequestHeader struct {
	Channel   string `header:"x-channel"`
	DeviceId  string `header:"x-device-id"`
	UserAgent string `header:"user-agent"`
	FcmToken  string `header:"x-fcm-token"`
}
