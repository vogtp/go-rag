import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';

@Injectable()
export class CollectionListService {

  collectionURL = '/search/';

  getCollections(): Observable<CollectionRequest> {
    console.log("Rest to "+this.collectionURL);
    
    return this.http.get<CollectionRequest>(this.collectionURL, { headers: httpHeaders }).pipe(
      catchError(this.handleError)
    );
  }


  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      console.log(error.error.message)

    } else {
      console.log(error.status)
    }
    return throwError(
      console.log('Something is wrong!'));
  };

  constructor(private http: HttpClient) { }
}

export interface CollectionRequest {
  Title: string;
  Version: string;
  Collections: Collection[];
}

export interface Collection {
  Name: string;
}

const httpHeaders: HttpHeaders = new HttpHeaders({
  //Authorization: 'Bearer JWT-token'
  Accept: "application/json"
})

