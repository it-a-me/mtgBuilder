import { Menubar } from "@/components/ui/menubar";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";
import react, { Suspense, useContext } from "react";
import { CardContext, CardsReady } from "./cardQuery";
import { Switch } from "./components/ui/switch";

function CardSearch({ setDisplayCards }: { setDisplayCards: react.Dispatch<react.SetStateAction<number[]>> }) {
  const searchRef = react.useRef<HTMLInputElement>(null)
  const [liveSearch, setLiveSearch] = react.useState(false)
  const cardQuery = useContext(CardContext)
  async function handleSubmit(e: react.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const query = searchRef.current?.value
    if (query == null) {
      alert("query null")
      return
    }
    const cards = await cardQuery.queryCards(query)
    setDisplayCards(cards)
  }
  return (
    <>
      <form className="flex w-full" onSubmit={e => handleSubmit(e)} onChange={e => liveSearch && handleSubmit(e)}>
        <Input className="bg-gray-300" ref={searchRef} type="search" placeholder="name:goblin type:creature oracle:/create.*token/" />
        <div className="flex items-center px-1">
          <p>Live</p>
          <Switch checked={liveSearch} onCheckedChange={() => setLiveSearch(v => !v)} className="flex" ></Switch>
        </div>
        <Button className="bg-gray-200" type="submit" variant="outline" >
          Search
        </Button>
      </form>
    </>
  );
}

function TitleBar({ setDisplayCards }: { setDisplayCards: react.Dispatch<react.SetStateAction<number[]>> }) {
  return (
    <>
      <Menubar className="pl-0">
        <Button className="bg-blue-300 hover:bg-blue-400 text-gray-800 font-bold" asChild>
          <a href="/">MtgBuilder</a>
        </Button>
        <CardSearch setDisplayCards={setDisplayCards} />
      </Menubar >
    </>
  );
}

function Card({ id, className }: { id: number, className?: string }) {
  const cardQuery = useContext(CardContext)
  const [card, setCard] = react.useState<{ imageUrl: string, cardUrl: string } | null>(null)
  react.useEffect(() => {
    let isMounted = true
    async function getCardUrl() {
      const card = await cardQuery.getCard(id)
      if (isMounted) {
        setCard({
          imageUrl: card.image_uris.border_crop.toString(),
          cardUrl: card.scryfall_uri.toString()
        })
      }
    }
    getCardUrl()
    return () => { isMounted = false }
  });
  return (<>
    {card != null && <a href={card.cardUrl} target="_blank">
      <img className={className} src={card.imageUrl} />
    </a>}
  </>)
}

function Body({ displayCards }: { displayCards: number[] }) {
  react.use(CardsReady)
  const MAX_CARDS = 30

  return <>{
    displayCards.length > 0 &&
    <h1 className="text-white p-4 font-bold">
      Displaying {Math.min(MAX_CARDS, displayCards.length)}/{displayCards.length}
    </h1>
  }
    < div className="flex flex-wrap justify-center" >
      {
        displayCards.slice(0, MAX_CARDS).map(c => <Card className="w-80 p-2" key={c} id={c} />)
      }
    </div ></>
}

function App() {
  const [displayCards, setDisplayCards] = react.useState<number[]>([])
  return (
    <>
      <div className="bg-gray-800 h-full min-h-screen">
        <div>
          <TitleBar setDisplayCards={setDisplayCards} />
        </div>
        <Suspense fallback=<h1 className="text-white text-2xl p-3">Loading Cards...</h1>>
          <Body displayCards={displayCards} />
        </Suspense>
      </div>
    </>
  );
}

export default App;
