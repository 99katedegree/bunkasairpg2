"use client";

export function LoadingScreen() {
  return (
    <div className="h-screen w-screen flex flex-col justify-center items-center text-4xl bg-[url('/bg-reward.png')] bg-[length:108%_108%] bg-no-repeat bg-center">
      <p className="font-dot-gothic">
        <span className="tracking-[8px]">Loading</span> ...
      </p>
      <div className="flex gap-5 mt-5">
        {[0, 1, 2, 3, 4, 5].map((i) => (
          <div
            key={i}
            className="w-7 h-10 bg-[url('/loading_slime.png')] bg-no-repeat bg-center bg-contain"
            style={{
              animation: "bounce 1.5s infinite",
              animationDelay: `${i * 0.2}s`,
            }}
          />
        ))}
      </div>
    </div>
  );
}
