'use client'
import styles from '@/app/styles/proposals.module.css'
import React, { useState } from "react"
import { LoremIpsum } from "lorem-ipsum";

export default function Proposals() {


  const [title, setTitle] = useState();
  const [author, setAuthor] = useState();
  const [email, setEmail] = useState();
  const [description, setDescription] = useState();
  const [generated, setGenerated] = useState(false);
  const [loading, setLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [sended, setSended] = useState(false);
  const [data, setData] = useState(null);


  const handleSubmit = () => {
    setLoading(true);
    setIsError(false);
    const data = {
      title: title,
      author: author,
      email: email,
      description: description,
    }

    console.log("Sending Post!" + JSON.stringify(data))
    try{
      fetch('/api/c4p', {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          'accept': 'application/json',
        },
      }).then((response) => response.json())
      .then((data) => {
        setData(data);
        setTitle('');
        setAuthor('');
        setEmail('');
        setDescription('');
        setLoading(false);
        setSended(true);
      })
    }catch(err){
        setLoading(false);
        setIsError(true);
      }
  }

  const lorem = new LoremIpsum({
    sentencesPerParagraph: {
      max: 8,
      min: 4
    },
    wordsPerSentence: {
      max: 16,
      min: 4
    }
  });

  
  function generate(){
    setDescription(lorem.generateParagraphs(2));
    setTitle(lorem.generateWords(5));
    setAuthor(lorem.generateWords(2));
    setEmail(lorem.generateWords(1)+"@mail.com");
    setGenerated(true);
  }

  const handleBack = () => {
    setSended(false)
  }

  return (
    <main className={styles.main}>
      <h1>Proposals</h1>

      <h2>Submit your Proposal</h2>

      <div>
      {!sended && (
        <div>
          <div>
            <label>Title</label>
            <input type="text" id="title" name="title" value={title}/>
          </div>

          <div>
            <label>Description</label>
            <textarea id="description" name="description" value={description}></textarea>
          </div>

          <div>
            <label>Author</label>
            <input type="text" id="author" name="author" value={author} />
          </div>

          <div>
            <label>Email</label>
            <input type="text" id="email" name="email" value={email} />
          </div>

          {isError && <small className="mt-3 d-inline-block text-danger">Something went wrong. Please try again later.</small>}

          {!generated && (  
              <button main inline onClick={generate} disabled={generated}>Generate</button>
          )}
          {generated && (
          <button type="submit" onClick={handleSubmit} >Send Proposal</button>
          )}
          </div>
          )}
          {sended && (
            <>
              <h3>Thanks!</h3>
              <button main onClick={handleBack} >Back</button>
            </>
          )}
      
    </div>
       
    </main>
  )
}
