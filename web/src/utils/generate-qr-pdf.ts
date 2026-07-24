import type { TDocumentDefinitions } from "pdfmake/interfaces";
import pdfMake from "pdfmake/build/pdfmake";
import * as pdfFonts from "pdfmake/build/vfs_fonts.js";
import QRCodeStyling from "qr-code-styling";

(pdfMake as any).vfs =
  (pdfFonts as any).default?.pdfMake?.vfs ||
  (pdfFonts as any).pdfMake?.vfs ||
  (pdfFonts as any).vfs;

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
  const qrImagePromises = urls.map((url) => generateTransparentQR(url));
  const qrDataUrls = await Promise.all(qrImagePromises);

  const content = qrDataUrls.map((dataUrl, index) => ({
    stack: [
      {
        image: dataUrl,
        width: 250,
        height: 250,
        alignment: "center" as const,
        margin: [0, 80, 0, 20] as [number, number, number, number],
      },
      {
        text: urls[index],
        fontSize: 8,
        alignment: "center" as const,
        color: "#666666",
        margin: [0, 0, 0, 0] as [number, number, number, number],
      },
    ],
    pageBreak: index === qrDataUrls.length - 1 ? undefined : ("after" as const),
  }));

  const docDefinition: TDocumentDefinitions = {
    pageSize: "A4",
    pageMargins: [40, 60, 40, 60],
    content,
  };

  pdfMake.createPdf(docDefinition).download(filename);
}
