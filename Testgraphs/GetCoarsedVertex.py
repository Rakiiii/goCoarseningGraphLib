import pandas as pd
import pickle

clf = pickle.load(open('model.pkl', 'rb'))

vertexPairs = pd.read_csv(r'pairs.csv')
predictions = clf.predict_proba(vertexPairs.values)
resFile = open('contractedPairs', 'w')
for i in range(len(predictions)):
    if predictions[i][1] >= 0.80:
        resFile.write("1 ")
    else:
        resFile.write("0 ")
resFile.close()
