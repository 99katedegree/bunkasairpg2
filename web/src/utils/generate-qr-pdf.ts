import type { TDocumentDefinitions } from "pdfmake/interfaces";
import pdfMake from "pdfmake/build/pdfmake";
import * as pdfFonts from "pdfmake/build/vfs_fonts.js";
import QRCodeStyling from "qr-code-styling";

(pdfMake as any).vfs =
  (pdfFonts as any).default?.pdfMake?.vfs ||
  (pdfFonts as any).pdfMake?.vfs ||
  (pdfFonts as any).vfs;

async function toDataUrl(url: string): Promise<string> {
  const res = await fetch(url);
  const blob = await res.blob();
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}

async function generateTransparentQR(url: string): Promise<string> {
  const qrCode = new QRCodeStyling({
    width: 250,
    height: 250,
    data: url,
    dotsOptions: { color: "#000000", type: "rounded" },
    backgroundOptions: { color: "transparent" },
  });

  const blob = (await qrCode.getRawData("png")) as Blob;
  if (!blob) throw new Error("QRコード生成に失敗しました");

  return new Promise<string>((resolve) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.readAsDataURL(blob);
  });
}

export async function generateQrPdf(urls: string[], filename: string): Promise<void> {
  const baseUrl = window.location.origin;
  const [logoDataUrl, bgDataUrl] = await Promise.all([
    toDataUrl(`${baseUrl}/logo.png`),
    toDataUrl(`${baseUrl}/bg-reward.png`),
  ]);

  const qrDataUrls = await Promise.all(urls.map((url) => generateTransparentQR(url)));

  const content = qrDataUrls.map((dataUrl, index) => ({
    stack: [
      {
        image: logoDataUrl,
        width: 480,
        alignment: "center" as const,
        margin: [0, 60, 0, 40] as [number, number, number, number],
      },
      {
        image: dataUrl,
        width: 250,
        height: 250,
        alignment: "center" as const,
        margin: [0, 0, 0, 20] as [number, number, number, number],
      },
    ],
    pageBreak: index === qrDataUrls.length - 1 ? undefined : ("after" as const),
  }));

  const docDefinition: TDocumentDefinitions = {
    pageSize: "A4",
    pageMargins: [40, 60, 40, 60],
    background: [
      {
        image: bgDataUrl,
        width: 580,
        height: 820,
        absolutePosition: { x: (595.28 - 580) / 2, y: (841.89 - 820) / 2 },
      },
    ],
    content,
  };

  pdfMake.createPdf(docDefinition).download(filename);
}
