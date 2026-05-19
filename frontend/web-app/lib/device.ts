"use client";

export type LocalDeviceSignal = {
  user_agent: string;
  viewport_width: number;
  viewport_height: number;
  screen_width: number;
  screen_height: number;
  has_touch: boolean;
  had_keyboard_event: boolean;
  had_pointer_event: boolean;
  correlation_id: string;
};

export function collectDeviceSignal(correlationId: string, interaction: { keyboard: boolean; pointer: boolean }): LocalDeviceSignal {
  return {
    user_agent: navigator.userAgent,
    viewport_width: window.innerWidth,
    viewport_height: window.innerHeight,
    screen_width: window.screen.width,
    screen_height: window.screen.height,
    has_touch: navigator.maxTouchPoints > 0,
    had_keyboard_event: interaction.keyboard,
    had_pointer_event: interaction.pointer,
    correlation_id: correlationId
  };
}

export function looksMobileLocally() {
  if (typeof navigator === "undefined") {
    return false;
  }
  const ua = navigator.userAgent.toLowerCase();
  return /mobile|android|iphone|ipad|tablet/.test(ua) || window.innerWidth < 1024;
}

