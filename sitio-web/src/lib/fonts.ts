import { Mozilla_Headline, Mozilla_Text } from "next/font/google";

export const mozillaText = Mozilla_Text({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-mozilla-text"
});

export const mozillaHeadline = Mozilla_Headline({
  subsets: ["latin"],
  display: "swap"
});
