import { Menubar } from "@/components/ui/menubar";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";
import react from "react";

function CardSearch({ setCard: setCardUrl }: { setCard: react.Dispatch<react.SetStateAction<null>> }) {
  const searchRef = react.useRef<HTMLInputElement>(null)
  function handleSubmit(e: react.FormEvent<HTMLFormElement>) {
    e.preventDefault();
  }
  return (
    <>
      <form className="flex w-full" onSubmit={e => handleSubmit(e)}>
        <Input ref={searchRef} type="search" placeholder="name:goblin type:creature oracle:/create.*token/" />
        <Button type="submit" variant="outline" >
          Search
        </Button>
      </form>
    </>
  );
}

function TitleBar({ setCard: setCardUrl }: { setCard: react.Dispatch<react.SetStateAction<null>> }) {
  return (
    <>
      <Menubar>
        <Button asChild>
          <a href="/">MtgBuilder</a>
        </Button>
        <CardSearch setCard={setCardUrl} />
      </Menubar >
    </>
  );
}

function App() {
  const [cardUrl, setCardUrl] = react.useState(null)
  return (
    <>
      <TitleBar setCard={setCardUrl} />
      {
        cardUrl != null && <img src={cardUrl} />
      }
    </>
  );
}

export default App;
