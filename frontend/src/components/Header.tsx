// Components
import { Button, ButtonType } from "./Button";
import { FileUpload } from "./FileUpload";

export const Header = () => {
  return (
    <header className="flex justify-end px-12">
      <FileUpload />
      <Button text="Login" type={ButtonType.LOGIN} />
    </header>
  );
};
