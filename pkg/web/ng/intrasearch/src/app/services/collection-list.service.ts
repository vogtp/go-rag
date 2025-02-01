import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { collectionURL, httpHeaders } from './common';
  
@Injectable()
export class CollectionListService {

  getCollections(): Observable<CollectionListResponse> {
    console.log("Rest to "+collectionURL);
    return this.http.get<CollectionListResponse>(collectionURL, { headers: httpHeaders }).pipe(
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

export interface CollectionListResponse {
  Title: string;
  Version: string;
  Collections: Collection[];
}

export interface Collection {
  Name: string;
}

