import { Link } from "react-router-dom";

interface IButtonProps {
  text: string;
  onClick?: () => void;
  url?: string;
  styles?: string | IButtonStyles;
}

export enum IButtonStyles {
  PRIMARY = "text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide",
}

export const Button = ({ text, onClick, url, styles }: IButtonProps) => {
  if (url) {
    return (
      <Link to={url} className={styles || IButtonStyles.PRIMARY}>
        {text}
      </Link>
    );
  }

  return (
    <button className={styles || IButtonStyles.PRIMARY} onClick={onClick}>
      {text}
    </button>
  );
};
