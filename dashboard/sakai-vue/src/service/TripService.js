import request from './wrapper';

export class TripService {
    createTrip(trip) {
        return request({
            method: 'post',
            url: '/trips',
            data: trip
        });
    }
    getAllTrip(){
        return request({
            method: 'get',
            url: '/trips'
        });
    }
    getThePageTrip(PageNumber){
        return request({
            method: 'get',
            url: '/trips?page=' + PageNumber
        });
    }
    getTrip(id){
        return request({
            method: 'get',
            url: '/trips/'+id
        });
    }
    joinTrip(id){
        return request({
            method: 'post',
            url: '/trips/'+id+'/join'
        });
    }
}
