import { Section, Heading, Img } from '@react-email/components';
import type { ReactElement, FC, PropsWithChildren } from 'react';

export default function EmailHeader({ children }: PropsWithChildren): ReactElement<FC> {
  return (
    <>
      <Section className="mt-8">
        <Img
          src="https://raw.githubusercontent.com/VEDA95/OpenBoard-Backend/refs/heads/main/docs/openboard_logo.png"
          width={64}
          height={64}
          alt="Open Board"
          className="mx-auto my-0"
        />
      </Section>
      <Heading className="mx-0 my-7 p-0 text-center font-normal text-[24px] text-black">
        {children}
      </Heading>
    </>
  );
}
