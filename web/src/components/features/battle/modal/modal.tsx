import Image from "next/image";

type Props = {
  onClose: () => void;
  onConfirm: () => void;
  title: string;
};

export function Modal({ onClose, onConfirm, title }: Props) {
  return (
    <div
      className="fixed top-0 left-0 w-screen h-screen flex items-center justify-center bg-black/80 bg-opacity-50 z-50"
      onClick={onClose}
    >
      <div
        className="flex flex-col justify-center items-center h-64 gap-8 p-10 rounded shadow-lg text-center bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: `url(${"/bg-reward.png"})` }}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mb-4 text-black whitespace-pre-wrap text-xl">{title}</div>
        <div className="flex gap-8 justify-center">
          <Image
            className="w-[130px] h-auto rounded"
            src={"/btn-yes.png"}
            alt="はい"
            width={400}
            height={400}
            onClick={onConfirm}
          />
          <Image
            className="w-[130px] h-auto rounded"
            src={"/btn-no.png"}
            alt="いいえ"
            width={400}
            height={400}
            onClick={onClose}
          />
        </div>
      </div>
    </div>
  );
}
