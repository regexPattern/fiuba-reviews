import { describe, expect, it } from "vitest";
import { detectarPlataforma } from "./plataforma";

const UA = {
  windows:
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
  linux:
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
  mac: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15",
  android:
    "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36",
  iphone:
    "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1"
};

describe("detectarPlataforma", () => {
  it("reconoce Windows y Linux como la misma plataforma", () => {
    expect(
      detectarPlataforma({ userAgent: UA.windows, platform: "Windows", maxTouchPoints: 0 })
    ).toBe("windows");
    expect(
      detectarPlataforma({ userAgent: UA.linux, platform: "Linux x86_64", maxTouchPoints: 0 })
    ).toBe("windows");
  });

  it("reconoce macOS, con platform viejo o de userAgentData", () => {
    expect(detectarPlataforma({ userAgent: UA.mac, platform: "MacIntel", maxTouchPoints: 0 })).toBe(
      "mac"
    );
    expect(detectarPlataforma({ userAgent: UA.mac, platform: "macOS", maxTouchPoints: 0 })).toBe(
      "mac"
    );
  });

  it("reconoce Android aunque platform diga Linux", () => {
    expect(
      detectarPlataforma({ userAgent: UA.android, platform: "Linux armv8l", maxTouchPoints: 5 })
    ).toBe("android");
  });

  it("reconoce el iPhone y el iPad que se hace pasar por Mac", () => {
    expect(
      detectarPlataforma({ userAgent: UA.iphone, platform: "iPhone", maxTouchPoints: 5 })
    ).toBe("ios");
    expect(detectarPlataforma({ userAgent: UA.mac, platform: "MacIntel", maxTouchPoints: 5 })).toBe(
      "ios"
    );
  });
});
