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
    "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1",
  ipad: "Mozilla/5.0 (iPad; CPU OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1"
};

describe("detectarPlataforma", () => {
  it("toma Windows y Linux como PC", () => {
    expect(detectarPlataforma(UA.windows)).toBe("pc");
    expect(detectarPlataforma(UA.linux)).toBe("pc");
  });

  it("reconoce macOS, con el iPad incluido", () => {
    expect(detectarPlataforma(UA.mac)).toBe("mac");
    expect(detectarPlataforma(UA.ipad)).toBe("mac");
  });

  it("reconoce Android aunque el UA diga Linux", () => {
    expect(detectarPlataforma(UA.android)).toBe("android");
  });

  it("reconoce el iPhone", () => {
    expect(detectarPlataforma(UA.iphone)).toBe("ios");
  });
});
