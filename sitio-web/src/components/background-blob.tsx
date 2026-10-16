export function BackgroundBlob({
  upperLeft = true,
  lowerRight = true
}: {
  upperLeft?: boolean;
  lowerRight?: boolean;
}) {
  return (
    <div
      className="pointer-events-none fixed inset-0 z-0 overflow-clip [contain:paint]"
      aria-hidden="true"
    >
      {upperLeft && (
        <div
          className="absolute -top-64 -left-64 h-[min(52rem,68vh)] w-[min(76rem,100vw)] bg-no-repeat max-sm:-top-48 max-sm:-left-52 max-sm:h-[48vh] max-sm:w-[125vw]"
          style={{
            backgroundImage:
              "radial-gradient(ellipse at center, rgb(79 172 212 / 0.2) 0%, rgb(79 172 212 / 0.12) 40%, rgb(79 172 212 / 0) 72%)"
          }}
        />
      )}
      {lowerRight && (
        <div
          className="absolute -right-40 -bottom-40 h-[min(26rem,42vh)] w-[min(30rem,65vw)] bg-no-repeat max-sm:-right-28 max-sm:-bottom-28 max-sm:h-[30vh] max-sm:w-[75vw]"
          style={{
            backgroundImage:
              "radial-gradient(ellipse at center, rgb(79 172 212 / 0.16) 0%, rgb(79 172 212 / 0) 70%)"
          }}
        />
      )}
    </div>
  );
}
