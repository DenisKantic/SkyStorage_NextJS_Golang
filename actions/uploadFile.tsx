"use server"

export async function uploadFile(formData:FormData){

    const response = await fetch("http://localhost:8080/upload",{
        method: "POST",
        body: formData,
        mode: "cors"
    })

    if(!response.ok){
        throw new Error("Network response is not okay")
    }

    const data = await response.json();
    

    return data;
}