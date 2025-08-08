import { Html, Head, Preview, Body, Container, Section, Text, Tailwind, pixelBasedPreset } from '@react-email/components';
import EmailHeader from '@components/header';
import type { ReactElement, FC } from 'react';

function PasswordResetEmailTemplate(): ReactElement<FC> {
  return (
    <Html>
      <Head />
      <Preview>Test...</Preview>
      <Tailwind config={{ presets: [pixelBasedPreset] }}>
        <Body className="mx-auto my-auto bg-white px-2 font-sans">
          <Container className="mx-auto my-10 max-w-[465px] rounded border border-[#eaeaea] p-5">
            <EmailHeader>
              <strong>Password Reset</strong>
            </EmailHeader>
            <Text className="text-4 text-black leading-6">
              Hello {'{{ .Username }}'},
            </Text>
            <Text className="text-4 text-black">
              A request was made recently for your account password to be reset. Please use the following OTP to initiate the password reset...
            </Text>
            <Section>
              <Text>{'{{ .Token }}'}</Text>
            </Section>
            <Text className="text-4 text-black">
              If you didn't request to have your password reset please disregard this email.
            </Text>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
}

export default PasswordResetEmailTemplate;
