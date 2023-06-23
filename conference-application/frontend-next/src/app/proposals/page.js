'use client'
import styles from '@/app/styles/proposals.module.css'
import React, { useState } from "react"
import { LoremIpsum } from "lorem-ipsum";
import Textfield from '../components/forms/textfield/textfield';
import Textarea from '../components/forms/textarea/textarea';
import Button from '../components/forms/button/button';
import toast, { Toaster } from "react-hot-toast";
import Cloud from '../components/cloud/cloud'

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
      fetch('/api/c4p/', {
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
      <div className={`${styles.hero} ` }>
        <div className={ `grid content noMargin`}>
          <div className="col full">
          <h1>Proposals <Cloud number="3" blue /></h1>
            
          </div>
        </div>
      </div>
    

      <div className="grid content">
        <div className="col third positionSingle">
        <h4>Join us as a speaker</h4>
        <p data-scroll data-scroll-speed="2" className="p p-b">Are you passionate about Cloud, Kubernetes, Docker or other technologies related with the Cloud. Submit your proposal to share your knowledge with our amazing community!</p>
         
        </div>
        <div className="col half positionHalf">
        {!sended && (
        <div>
          
            <Textfield label="Title" id="title" name="title" value={title} />
            <Textarea label="Description" id="description" name="description" value={description}  />  
            
            <Textfield label="Author" id="author" name="author" value={author} />
            <Textfield label="Email" id="email" name="email" value={author} />
            

          {isError && <small className="mt-3 d-inline-block text-danger">Something went wrong. Please try again later.</small>}

          {!generated && (  
              <Button main clickHandler={generate} disabled={generated}>Generate</Button>
          )}
          {generated && (
          <Button type="submit" clickHandler={handleSubmit} >Send Proposal</Button>
          )}
          </div>
          )}
          {sended && (
            <>
              <h3>Thanks!</h3> 
              <Button  clickHandler={handleBack} >Send another proposal</Button>
            </>
          )}
        </div>
      </div>
      

      

      <div>
      
      
    </div>
       
    </main>
  )
}
