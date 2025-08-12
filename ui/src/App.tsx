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
        <Input className="bg-gray-300 text-xl" ref={searchRef} type="search" placeholder="name:goblin type:creature oracle:/create.*token/" />
        <div className="flex items-center px-1">
          <p className="text-xl">Live</p>
          <Switch checked={liveSearch} onCheckedChange={() => setLiveSearch(v => !v)} className="flex" ></Switch>
        </div>
        <Button className="bg-gray-200 text-xl" type="submit" variant="outline" >
          Search
        </Button>
      </form>
    </>
  );
}

function TitleBar({ setDisplayCards }: { setDisplayCards: react.Dispatch<react.SetStateAction<number[]>> }) {
  return (
    <>
      <Menubar className="p-0">
        <Button className="bg-blue-300 hover:bg-blue-400 text-xl text-gray-800 font-bold" asChild>
          <a href="/">MtgBuilder</a>
        </Button>
        <CardSearch setDisplayCards={setDisplayCards} />
      </Menubar >
    </>
  );
}

function Card({ getCard: cardPromise, className }: { getCard: Promise<object>, className?: string }) {
  const card = react.use(cardPromise)

  let name, imageUrl, cardUrl;
  try {
    //@ts-expect-error Untyped value
    name = card.name
    //@ts-expect-error Untyped value
    cardUrl = card.scryfall_uri
    //@ts-expect-error Untyped value
    imageUrl = card.image_uris.border_crop
  }
  catch (e) {
    if (e instanceof TypeError) {
      return <>
        <a href={cardUrl} className={className} title={name}>
          <h1 className="bg-white text-red-500 text-2xl h-full px-4 py-2">
            {name} has no associated image
          </h1>
        </a>
      </>
    }
    throw e
  }

  return (<>
    {card != null && <a href={cardUrl} target="_blank">
      <img className={className} src={imageUrl} title={name} />
    </a>}
  </>)
}

function Body({ displayCards }: { displayCards: number[] }) {
  const MAX_CARDS = 60
  react.use(CardsReady)
  const cardQuery = useContext(CardContext)
  const [page, setPage] = react.useState<number>(0)
  react.useEffect(() => setPage(0), [displayCards])
  const startIndex = page * MAX_CARDS
  const endIndex = Math.min(page * MAX_CARDS + MAX_CARDS, displayCards.length)
  const lastPage = Math.max(0, Math.trunc(displayCards.length / MAX_CARDS))

  return <>
    <div className="flex items-center justify-between p-1">
      <h1 className="text-white text-xl font-bold">
        Displaying {startIndex}-{endIndex} of {displayCards.length} matching cards
      </h1>
      <div>
        <Button className="mx-1" variant="outline" onClick={() => setPage(p => Math.max(p - 1, 0))}>
          Previous
        </Button>
        <Button className="mx-1" variant="outline" onClick={() => setPage(p => Math.min(p + 1, lastPage))}>
          Next
        </Button>
      </div>
    </div>
    <div className="flex flex-wrap justify-center" >
      {
        displayCards.slice(startIndex, startIndex + MAX_CARDS).map(id =>
          <Suspense key={id}>
            <Card className="w-100 p-2" getCard={cardQuery.getCard(id)} />
          </Suspense>)
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
