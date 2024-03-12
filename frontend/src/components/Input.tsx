interface InputProps {
  placeholder: string;
  styles: string;
}

export const Input = ({ placeholder, styles }: InputProps) => {
  return <input className={styles} placeholder={placeholder} />;
};
