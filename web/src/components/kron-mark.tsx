type Props = { className?: string };

export function KronMark({ className }: Props) {
  return (
    <svg
      viewBox="0 0 32 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden
    >
      <circle cx="16" cy="16" r="14" stroke="currentColor" strokeWidth="1.25" />
      <circle cx="16" cy="16" r="9" stroke="currentColor" strokeWidth="1" opacity="0.55" />
      <circle cx="16" cy="16" r="4.5" stroke="currentColor" strokeWidth="1" opacity="0.35" />
      <line x1="16" y1="16" x2="16" y2="6.5" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
      <line x1="16" y1="16" x2="22.5" y2="19" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
      <circle cx="16" cy="16" r="1.4" fill="currentColor" />
    </svg>
  );
}
