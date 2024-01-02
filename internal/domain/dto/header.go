package dto

type RequestHeader struct {
	Authorization string `header:"Authorization"`
	Channel       string `header:"x-channel"`
	DeviceId      string `header:"x-device-id" validate:"required"`
	UserAgent     string `header:"user-agent"`
	FcmToken      string `header:"x-fcm-token"`
}
