from cmath import e
import numpy as np
import matplotlib.pyplot as plt
import scipy.stats as st
import seaborn as sns
import pandas as pd
import matplotlib.ticker as tic

sns.set_theme(style="darkgrid")

folder = "data/"

filename = "q,k=2,lam=100,D=100,orphanage"
xlims = [0.4, 1]
xlabel = "Adversary proportion"

# values
lam = 100.
h = 1.
k = 2.
D = 100.

printAnalytical = True
ylimsTips = [0, 10000]
ylimsOrphan = [1e-6, 1]
folderdata = "../data/"+filename+"/"

# Colors
BG_WHITE = "#fbf9f4"
GREY_LIGHT = "#b4aea9"
GREY50 = "#7F7F7F"
BLUE_DARK = "#1B2838"
BLUE = "#2a475e"
BLACK = "#282724"
GREY_DARK = "#747473"
RED_DARK = "#850e00"
# Colors taken from Dark2 palette in RColorBrewer R library
COLOR_SCALE = ["#1B9E77", "#D95F02", "#7570B3"]


def main():
    evaluate1(1)  # tips
    evaluate1(2)  # orphanage


def evaluate1(analysisType):
    if analysisType == 1:
        filenamedata = folderdata+"tips_"
        ylabel = "Number of tips"
        fileSaveFig = folder+'tips'
        column = 1
    elif analysisType == 2:
        filenamedata = folderdata+"orphantips_"
        ylabel = "Orphanage rate"
        fileSaveFig = folder+'orphanage'
        column = 2

    X = loadColumn(folderdata+"params", 0, 0)
    Y = X*0.
    yQ1 = X*0.
    yQ3 = X*0.
    yMin = X*0.
    yMax = X*0.

    fig, ax = plt.subplots()
    for i in np.arange(len(X)):
        y_data = loadColumn(filenamedata+str(i), column, 2)
        # only consider data when it is converged
        y_data = y_data[int(len(y_data)*.75):]
        Y[i] = np.mean(y_data)
        df = pd.DataFrame(y_data, columns=['data'])
        dfStats = df['data'].describe()
        yQ1[i] = dfStats['25%']
        yQ3[i] = dfStats['75%']
        # 0 values make no sense
        Y[Y == 0] = np.nan
        yQ3[yQ1 == 0] = np.nan
        yQ1[yQ1 == 0] = np.nan
        yMax[yMin == 0] = np.nan
        yMin[yMin == 0] = np.nan
    sns.lineplot(x=X, y=Y, label="Simulation")
    # plt.fill_between(X, yQ1, yQ3, color='b',
    #                  alpha=0.2, label="25% to 75% quantiles")
    xL, L1, L2, Lavg, qSimple, Lsimple, o, oSimple = getAnalyticalCurve(
        X, Y, analysisType)
    if printAnalytical & (analysisType == 1):
        sns.lineplot(x=qSimple, y=Lsimple, color="red",
                     label="Analytical (simple)", linestyle="dashed")
        sns.lineplot(x=xL, y=Lavg, color="red", label="Analytical")
        sns.lineplot(x=xL, y=L1, color="red", linestyle="dotted")
        sns.lineplot(x=xL, y=L2, color="red", linestyle="dotted")
    if printAnalytical & (analysisType == 2):
        sns.lineplot(x=qSimple, y=oSimple, color="red",
                     label="Model A", linestyle="dashed")
        sns.lineplot(x=xL, y=o, color="purple",
                     label="Model B", linestyle="dotted")
    plt.ylim(ylimsTips)
    if analysisType == 2:
        plt.yscale('log')
        plt.ylim(ylimsOrphan)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.xlim(xlims)
    plt.legend()
    plt.savefig(fileSaveFig+'.png', format='png')
    plt.clf()

    # second type chart
    if analysisType == 1:
        # fig = plt.figure()
        # ax1=fig.add_subplot(211)
        # ax2=fig.add_subplot(212)
        setticks = [0, .25, .5, .75, 1.]
        fig, axes = plt.subplots(nrows=2, sharex=True, sharey=True)
        ax1 = plt.subplot(211)
        ax2 = plt.subplot(212)
        fig.subplots_adjust(hspace=.02)
        splity = 3000
        ax1.set_ylim([splity*1.01, 10000])
        ax2.set_ylim([0, splity])
        ax1.set_xlim([0, 1])
        ax2.set_xlim([0, 1])
        ax1.xaxis.set_ticklabels([])
        ax1.get_xaxis().set_ticks(setticks)
        ax2.get_xaxis().set_ticks(setticks)
        ax1.plot(X, Y, label="Simulation", color="tab:blue")
        ax2.plot(X, Y, color="tab:blue")
        cutIndexData = 0
        ax1.plot(qSimple[cutIndexData:], Lsimple[cutIndexData:], color="red",
                 label="Analytical (Model A)", linestyle="dashed")
        ax2.plot(qSimple, Lsimple, color="red",
                 label="Analytical (Model A)", linestyle="dashed")
        cutIndex2 = np.sum(xL < .6)
        ax1.plot(xL[:cutIndex2], Lavg[:cutIndex2],
                 color="purple", label="Analytical (Model B)", linestyle="dotted")
        ax2.plot(xL[:cutIndex2], Lavg[:cutIndex2],
                 color="purple", label="Analytical (Model B)", linestyle="dotted")
        ax1.legend()
        ax1.set_xlabel("Tip pool size")
        ax2.set_ylabel("Adversary proportion")
        plt.savefig(fileSaveFig+'_subplots.png', format='png')
        plt.clf()


def getAnalyticalCurve(Xdata, Ydata, analysistype):
    # load X data
    q = (np.arange(1000)/1000.)*.5*.999
    L1, L2, L0avg, delta = calcL(q)

    # simple equation
    qSimple = q
    Lsimple = (1-qSimple)*k/((1-qSimple)*k-1)*h*lam*np.ones(len(qSimple))

    # get it from file instead
    q, L0avg = getLfromWolfram("data/dataWolfram/k=2,D=100,lam=100")
    L1 = L0avg*0  # delete because of Wolfram
    L2 = L0avg*0  # delete because of Wolfram

    Lsimple[Lsimple < 0] = np.NaN

    # estimate coefficient for orphanage probability
    if analysistype == 2:
        # need to ignore data for which probability is 0
        # Xdata = Xdata[Ydata > 0]
        # Ydata = Ydata[Ydata > 0]
        # _, _, L0avgdata, _ = calcL(Xdata)
        # print("L0avgdata", len(L0avgdata), L0avgdata)
        # print("Xdata", len(Xdata), Xdata)
        # print("Ydata", len(Ydata), Ydata)
        # select = (Xdata >= 2.)*(Ydata < .1)
        # cOrph0 = Ydata[select] * np.exp(Xdata[select]*D*k / L0avgdata[select])
        # print("median cOrph0= ", np.median(np.sort(cOrph0)))
        # c = np.median(np.sort(cOrph0))
        print("-------------------------------")
        print("getAnalyticalCurve: for now set c=1")
        c = 1.
    else:
        c = 1.

    # c = 50.
    o = c*np.exp(-D*(1-q)*k/L0avg)
    oSimple = c*np.exp(-D*(1-qSimple)*k/(Lsimple/lam))

    return q, L1, L2, L0avg*lam, qSimple, Lsimple, o, oSimple


def calcL(q):
    # solve the following numerically
    # L0 = p*k(h-D*np.exp(-D*p*k/L0))/(p*k-1+np.exp(-D*p*k/L0))*lam*h
    p = 1-q
    L0 = q*0.+3.
    print("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
    print(h, k, D)
    # print(p)
    for i in range(1001):
        L1 = L0
        L0 = p*k*(h-D*np.exp(-D*p*k/L0))/(p*k-1+np.exp(-D*p*k/L0))
        delta = L0-L1
    L2 = L0*lam
    L1 = L1*lam
    # L1[L1 < 0] = np.NaN
    # L2[L1 < 0] = np.NaN
    L0avg = (L1+L2)/2./lam
    return L1, L2, L0avg, delta


def getLfromWolfram(file):
    x = loadColumn(file, 0, 2)
    y = loadColumn(file, 1, 2)
    return x, y/lam


def loadColumn(filename, column, skiprows):
    try:
        filestr = filename+".csv"
        f = open(filestr, "r")
        data = np.loadtxt(f, delimiter=";",
                          skiprows=skiprows, usecols=(column))
        return data
    except FileNotFoundError:
        print(filestr)
        print("File not found.")
        return []


# needs to be at the very end of the file
if __name__ == '__main__':
    main()
